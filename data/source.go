package data

// Source источник данных для сервиса
type Source interface {
	GetCities() ([]string, error)
}
