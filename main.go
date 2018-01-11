package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/matperez/city-autocomplete/persistence"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/mattn/go-sqlite3"
)

// создание таблицы для хранения данных
func ensureDB(db *sql.DB) error {
	sqlStmt := `
	create table foo (
		id integer not null primary key,
		name text,
		popularity integer not null default 0
	);
	delete from foo;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	logger := log.NewJSONLogger(os.Stderr)

	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}
	defer db.Close()

	store := persistence.NewSQLiteStore(db)

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
