package persistence

import (
	"database/sql"

	"github.com/go-kit/kit/log"
)

// реализация интерфейса через SQLite
type sqlStore struct {
	db     *sql.DB
	logger log.Logger
}

// NewSQLStore фабрика нового хранилища
func NewSQLStore(db *sql.DB, logger log.Logger) (Store, error) {
	store := &sqlStore{db, logger}
	if err := store.init(); err != nil {
		return store, err
	}
	return store, nil
}

// Init инициализация хранилища
func (s sqlStore) init() error {
	s.logger.Log("event", "Initializing the database")
	sqlStmt := `
	create table city (
		name text
	);
	delete from city;
	`
	_, err := s.db.Exec(sqlStmt)
	if err != nil {
		s.logger.Log("err", "Initialization failed")
		return err
	}
	s.logger.Log("event", "Database initialized")
	return nil
}

// Query запросить список удовлетворяющий условиям поиска
// query - строка поиска
// serviceType - тип услуги
func (s sqlStore) Query(query string) ([]string, error) {
	stmt, err := s.db.Prepare("select name from city where name like ?")
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	items := []string{}
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
func (s sqlStore) Populate(items []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("delete from city")
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("insert into city(name) values(?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(item)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
