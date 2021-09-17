package mysql

import "database/sql"

type Store struct {
	DB *sql.DB
}
