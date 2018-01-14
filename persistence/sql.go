package persistence

import (
	"database/sql"
	"strings"

	"github.com/essentialkaos/translit"

	"github.com/matperez/city-autocomplete/data"
)

// реализация интерфейса через SQLite
type sqlStore struct {
	db *sql.DB
}

// NewSQLStore фабрика нового хранилища
func NewSQLStore(db *sql.DB) (Store, error) {
	store := &sqlStore{db}
	if err := store.init(); err != nil {
		return store, err
	}
	return store, nil
}

// Init инициализация хранилища
func (s sqlStore) init() error {
	sqlStmt := `
	create table city (
		name string,
		code string,
		country_name string,
		country_code string,
		total_score integer,
		arrival_score integer,
		departure_score integer,
		transit_score integer,
		autocomplete string
	);
	delete from city;
	`
	_, err := s.db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

// Query запросить список удовлетворяющий условиям поиска
// query - строка поиска
// serviceType - тип услуги
func (s sqlStore) Query(query string) ([]data.City, error) {
	items := []data.City{}
	if query == "" {
		return items, nil
	}
	sqlQuery := `
	select 
		name, 
		code, 
		country_name, 
		country_code, 
		total_score, 
		arrival_score, 
		departure_score, 
		transit_score 
	from city where autocomplete like ? ORDER BY total_score DESC limit 10
	`
	stmt, err := s.db.Prepare(sqlQuery)
	if err != nil {
		return items, err
	}
	defer stmt.Close()
	rows, err := stmt.Query("%" + strings.ToLower(query) + "%")
	if err != nil {
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		var item data.City
		err = rows.Scan(&item.Name, &item.Code, &item.CountryName, &item.CountryCode, &item.TotalScore, &item.ArrivalScore, &item.DepartureScore, &item.TransitScore)
		if err != nil {
			return []data.City{}, nil
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, err
	}
	return items, nil
}

// Populate очистить базу данных и заполнить ее заново переданными значениями
func (s sqlStore) Populate(items []data.City) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("delete from city")
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStmt := `
	insert into city(
		name,
		code,
		country_name,
		country_code,
		total_score,
		arrival_score,
		departure_score,
		transit_score,
		autocomplete
	)
	values(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(
			item.Name,
			item.Code,
			item.CountryName,
			item.CountryCode,
			item.TotalScore,
			item.ArrivalScore,
			item.DepartureScore,
			item.TransitScore,
			autocomplete(item),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// подготавливает строку по которой будет выполняться поиск
func autocomplete(c data.City) string {
	return strings.ToLower(strings.Join([]string{
		c.Name,
		c.Code,
		translit.EncodeToISO9B(c.Name),
	}, " "))
}
