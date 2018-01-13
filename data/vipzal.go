package data

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type client struct {
	baseURL string
}

// NewVipzalSource фабричный метод
func NewVipzalSource(baseURL string) Source {
	return &client{baseURL}
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

	log.WithFields(log.Fields{
		"baseUrl": c.baseURL,
	}).Warn("Fetching cities")

	err := fetchJSON(c.baseURL+"/white-label/city/scored", target)

	if err != nil {
		return []City{}, err
	}

	return *target, nil
}
