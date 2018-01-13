package main

import (
	"time"

	"github.com/matperez/city-autocomplete/data"
	log "github.com/sirupsen/logrus"
)

// структура декоратора
type loggingMiddleware struct {
	next Service
}

// декоратор для логирования данных запроса
func (mw loggingMiddleware) Query(query string, serviceType string) (output []data.City, err error) {
	defer func(begin time.Time) {
		log.WithFields(log.Fields{
			"method": "query",
			"input":  map[string]string{"query": query, "type": serviceType},
			"output": output,
			"err":    err,
			"took":   time.Since(begin),
		}).Info("Fetching cities")
	}(time.Now())

	output, err = mw.next.Query(query, serviceType)
	return
}
