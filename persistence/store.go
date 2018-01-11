package persistence

// Store интерфейс хранилища данных
type Store interface {
	Populate([]string) error
	Query(string) ([]string, error)
}
