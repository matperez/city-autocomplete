package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/matperez/city-autocomplete/data"
	"github.com/matperez/city-autocomplete/persistence"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/mattn/go-sqlite3"
)

var (
	baseURL string
)

func init() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s [OPTIONS] <base-url>\n", os.Args[0])
		fmt.Println("Options available:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	baseURL = flag.Arg(0)
}

func main() {
	logger := log.NewJSONLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	os.Remove("./db.sqlite")

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	store, err := persistence.NewSQLStore(db, logger)
	if err != nil {
		panic(err.Error())
	}

	logger.Log("event", "Populating the database")

	vipzal := data.NewVipzalSource(baseURL, logger)
	cities, err := vipzal.GetCities()
	if err != nil {
		panic(err.Error())
	}

	if err = store.Populate(cities); err != nil {
		panic(err.Error())
	}

	logger.Log("event", "Database populated")

	var svc Service
	svc = New(store)
	svc = loggingMiddleware{logger, svc}

	http.Handle(
		"/query",
		httptransport.NewServer(
			makeQueryEndpoint(svc),
			decodeQueryRequest,
			encodeResponse,
		),
	)

	logger.Log("event", "Listening", "proto", "HTTP", "addr", ":8080")
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
	var query, serviceType string
	query = r.URL.Query().Get("query")
	serviceType = r.URL.Query().Get("type")
	return queryRequest{query, serviceType}, nil
}

// функция для кодирования ответа
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// обработчик запроса на получение списка городов
func makeQueryEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(queryRequest)
		v, err := svc.Query(req.Q, req.T)
		if err != nil {
			return queryResponse{v, err.Error()}, nil
		}
		return queryResponse{v, ""}, nil
	}
}
