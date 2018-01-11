package persistence

import (
	"database/sql"
)

// реализация интерфейса через SQLite
type sqliteStore struct {
	db *sql.DB
}

// NewSQLiteStore фабрика нового хранилища
func NewSQLiteStore(db *sql.DB) Store {
	return &sqliteStore{db}
}

// Query запросить список удовлетворяющий условиям поиска
// query - строка поиска
// serviceType - тип услуги
func (s sqliteStore) Query(query string) ([]string, error) {
	stmt, err := s.db.Prepare("select name from city where name like ?")
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return []string{}, nil
		}
		items = append(items, name)
	}
	err = rows.Err()
	if err != nil {
		return items, err
	}
	return items, nil
}

// Populate очистить базу данных и заполнить ее заново переданными значениями
func (s sqliteStore) Populate(items []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("delete from city")
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into city(name) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(item)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}
