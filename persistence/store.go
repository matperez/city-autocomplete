package persistence

import "github.com/matperez/city-autocomplete/data"

// Store интерфейс хранилища данных
type Store interface {
	Populate([]data.City) error
	Query(string) ([]data.City, error)
}
