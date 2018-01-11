package main

import "github.com/matperez/city-autocomplete/persistence"

// Service интерфейс сервиса для выполнения запросов
type Service interface {
	Query(string, string) ([]string, error)
}

// реализация сервиса
type service struct {
	store persistence.Store
}

// New фабричный метод
func New(store persistence.Store) Service {
	return &service{store}
}

func (s service) Query(query string, serviceType string) ([]string, error) {
	return s.store.Query(query)
}
