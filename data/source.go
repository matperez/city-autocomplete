package data

// City данные города
type City struct {
	Name           string `json:"city_name"`
	Code           string `json:"city_code"`
	CountryName    string `json:"country_name"`
	CountryCode    string `json:"country_code"`
	TotalScore     string `json:"total_score"`
	ArrivalScore   string `json:"arrival_score"`
	DepartureScore string `json:"departure_score"`
	TransitScore   string `json:"transit_score"`
}

// Source источник данных для сервиса
type Source interface {
	GetCities() ([]City, error)
}
