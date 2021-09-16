package store

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	DB *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Init() error {
	var err error
	s.DB, err = sql.Open("mysql", "root:password@(localhost:7200)/db?parseTime=true")
	if err != nil {
		return err
	}
	return nil
}
