package main

// Service интерфейс сервиса для выполнения запросов
type Service interface {
	Query(string, string) ([]string, error)
}

// реализация сервиса
type service struct{}

func (s service) Query(query string, serviceType string) ([]string, error) {
	return []string{query, serviceType}, nil
}
