package data

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type client struct {
	baseURL string
	logger  log.Logger
}

// NewVipzalSource фабричный метод
func NewVipzalSource(baseURL string, logger log.Logger) Source {
	return &client{baseURL, logger}
}

func fetchJSON(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(target)
	return nil
}

// Получение списка городов
func (c client) GetCities() ([]City, error) {
	target := new([]City)

	c.logger.Log("event", "Fetching cities", "baseUrl", c.baseURL)
	err := fetchJSON(c.baseURL+"/white-label/city/scored", target)

	if err != nil {
		return []City{}, err
	}

	return *target, nil
}
