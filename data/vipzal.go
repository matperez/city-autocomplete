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

// City структура для отображения данных города
type city struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	CountryName  string `json:"country_name"`
	Autocomplete string `json:"autocomplete"`
}

// Response ответ сервера
type response []city

// Получение списка городов
func (c client) GetCities() ([]string, error) {
	target := new(response)

	c.logger.Log("event", "Fetching cities", "baseUrl", c.baseURL)
	err := fetchJSON(c.baseURL+"/white-label/city", target)

	if err != nil {
		return []string{}, err
	}

	items := []string{}

	for _, city := range *target {
		items = append(items, city.Name)
	}

	c.logger.Log("event", "Cities fetched")
	return items, nil
}
