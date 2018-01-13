package main

import (
	"time"

	"github.com/matperez/city-autocomplete/data"

	"github.com/go-kit/kit/log"
)

// структура декоратора
type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

// декоратор для логирования данных запроса
func (mw loggingMiddleware) Query(query string, serviceType string) (output []data.City, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "query",
			"input", map[string]string{"query": query, "type": serviceType},
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Query(query, serviceType)
	return
}
