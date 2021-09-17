package mysql

import (
	"context"
	"errors"
	"url-shortner/domain"
)

func (s *Store) CreateUsersTable(ctx context.Context) error {
	table := `CREATE TABLE IF NOT EXISTS users (
    id BINARY(16) DEFAULT (UUID_TO_BIN(UUID())),
    email    VARCHAR(26) NOT NULL UNIQUE,
    passwordHash VARCHAR(26) NOT NULL,
    createdAt DATETIME DEFAULT NOW(),
    updatedAt DATETIME DEFAULT NOW(),
    PRIMARY KEY (id)
  );`
	_, err := s.DB.ExecContext(ctx, table)
	return err
}

func (s *Store) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	st, err := s.DB.PrepareContext(ctx, "INSERT INTO users(email,passwordHash) values(?,?) RETURNING id,createdAt,updatedAt")
	if err != nil {
		return nil, errors.New("could not prepare statment")
	}
	err = st.QueryRowContext(ctx, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, errors.New("could not insert row into users table")
	}
	return u, nil
}
