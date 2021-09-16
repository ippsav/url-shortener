package repository

import (
	"main/models"
	"main/store"
)

var (
	CREATE_USER_TABLE = `CREATE TABLE users (
    id VARCHAR(36),
    email    VARCHAR(26) NOT NULL,
    passwordHash VARCHAR(26) NOT NULL,
    createdAt DATETIME,
    updatedAt DATETIME,
    PRIMARY KEY (id),
    UNIQUE (email)
  );`
	INSERT_USER = `INSERT INTO users (id,email,passwordHash,createdAt) values(?,?,?,?)`
	QUERY_USER  = `SELECT id, email, passwordHash, createdAt FROM users WHERE email= ?`
	dbStore     = store.NewStore()
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) CreateUserTable() error {
	_, err := dbStore.DB.Exec(CREATE_USER_TABLE)
	return err
}

func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	_, err := dbStore.DB.Exec(INSERT_USER, u.ID, u.Email, u.PasswordHash, u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) FindUser(u *models.User) (*models.User, error) {
	err := dbStore.DB.QueryRow(QUERY_USER, u.Email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
