package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := log.NewJSONLogger(os.Stderr)

	var svc Service
	svc = service{}
	svc = loggingMiddleware{logger, svc}

	http.Handle(
		"/query",
		httptransport.NewServer(
			makeUppercaseEndpoint(svc),
			decodeQueryRequest,
			encodeResponse,
		),
	)

	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}

// структура для парсинга запроса
type queryRequest struct {
	Q string `json:"query"`
	T string `json:"type"`
}

// структура для кодирования ответа
type queryResponse struct {
	V   []string `json:"v"`
	Err string   `json:"err,omitempty"`
}

// функция для парсинга запроса
func decodeQueryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request queryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// функция для кодирования ответа
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// обработчик запроса на получение списка городов
func makeUppercaseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(queryRequest)
		v, err := svc.Query(req.Q, req.T)
		if err != nil {
			return queryResponse{v, err.Error()}, nil
		}
		return queryResponse{v, ""}, nil
	}
}
