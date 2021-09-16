package models

import (
	"database/sql"
	"fmt"
	"main/utils"
	"time"

	"github.com/google/uuid"
)

var (
	CREATE_USER_TABLE = `CREATE TABLE users (
    id VARCHAR(36),
    email    VARCHAR(20) NOT NULL,
    passwordHash VARCHAR(26) NOT NULL,
    createdAt DATETIME,
    updatedAt DATETIME,
    PRIMARY KEY (id),
    UNIQUE (email)
  );`
	INSERT_USER = `INSERT INTO users (id,email,passwordHash,createAt) values(?,?,?,?)`
	QUERY_USER  = `SELECT id, email, passwordHash, createdAt FROM users WHERE email= ?`
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	passwordHash string
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewUser() *User {
	return &User{}
}

func DBMigrate(db *sql.DB) error {
	_, err := db.Exec(CREATE_USER_TABLE)
	return err
}

func (u *User) CreateUser(db *sql.DB, email string, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("Invalid value for email or password\n")
	}
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u.ID = uuid.NewString()
	u.Email = email
	u.passwordHash = passwordHash
	u.CreatedAt = time.Now()
	_, err = db.Exec(INSERT_USER, u.ID, u.Email, u.passwordHash, u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
