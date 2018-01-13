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
	log "github.com/sirupsen/logrus"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/mattn/go-sqlite3"
)

var (
	baseURL  string
	logLevel string
)

func init() {
	flag.StringVar(&logLevel, "logLevel", "info", "Log level: warning, info, error, fatal")
	flag.Parse()
	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}
	baseURL = flag.Arg(0)
	setupLogger()
}

// наполнение базы данными с випзала
func populateDatabase(store persistence.Store) error {
	log.Warn("Populating the database")

	vipzal := data.NewVipzalSource(baseURL)
	cities, err := vipzal.GetCities()
	if err != nil {
		return err
	}

	if err = store.Populate(cities); err != nil {
		return err
	}

	log.Warn("Database populated")
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()

	store, err := persistence.NewSQLStore(db)
	if err != nil {
		log.Panic(err.Error())
	}

	if err = populateDatabase(store); err != nil {
		log.Panic(err.Error())
	}

	var svc Service
	svc = NewService(store)
	svc = loggingMiddleware{svc}

	http.Handle(
		"/query",
		httptransport.NewServer(
			makeQueryEndpoint(svc),
			decodeQueryRequest,
			responseJSON,
		),
	)

	http.Handle("/update", httptransport.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			log.Warn("Database repopulation requested")
			if err := populateDatabase(store); err != nil {
				return nil, err
			}
			return "ok", nil
		},
		func(_ context.Context, r *http.Request) (interface{}, error) {
			return nil, nil
		},
		responseString,
	))

	log.WithFields(log.Fields{
		"addr": ":8080",
	}).Warn("Listening")

	log.Error(http.ListenAndServe(":8080", nil))
}

// Настройка логирования
func setupLogger() {
	// Only log the warning severity or above.
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err.Error())
		printUsage()
		os.Exit(1)
	}
	log.SetLevel(level)
}

// отображение справки по параметрам запуска
func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] <base-url>\n", os.Args[0])
	fmt.Println("Options available:")
	flag.PrintDefaults()
}

// структура для парсинга запроса
type queryRequest struct {
	Q string `json:"query"`
	T string `json:"type"`
}

// структура для кодирования ответа
type queryResponse struct {
	V   []data.City `json:"v"`
	Err string      `json:"err,omitempty"`
}

// функция для парсинга запроса
func decodeQueryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var query, serviceType string
	query = r.URL.Query().Get("query")
	serviceType = r.URL.Query().Get("type")
	return queryRequest{query, serviceType}, nil
}

// функция для кодирования ответа в JSON
func responseJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// функция для кодирования ответа в текст
func responseString(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Fprint(w, response)
	return nil
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
