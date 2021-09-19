package mysql

import (
	"context"
	"url-shortner/domain"

	"github.com/pkg/errors"
)

func (s *Store) CreateUsersTable(ctx context.Context) error {
	table := `CREATE TABLE IF NOT EXISTS users (
    id BINARY(16),
    email    VARCHAR(26) NOT NULL UNIQUE,
    passwordHash VARCHAR(60) NOT NULL,
    createdAt DATETIME ,
    updatedAt DATETIME ,
    PRIMARY KEY (id)
  );`
	_, err := s.DB.ExecContext(ctx, table)
	return err
}

//SELECT BIN_TO_UUID(id),createdAt,updatedAt FROM users where email='?'
func (s *Store) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	st, err := s.DB.PrepareContext(ctx, "INSERT INTO users(id,email,passwordHash,createdAt,updatedAt) values(UUID_TO_BIN(?),?,?,?,?);")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	_, err = st.ExecContext(ctx, u.ID, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert row into users table")
	}
	// getting the user data back since mysql doesn t support returning keyword
	// st, err = s.DB.PrepareContext(ctx, "SELECT BIN_TO_UUID(id),createdAt,updatedAt FROM users where email=?")
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not prepare select statement")
	// }
	// err = st.QueryRowContext(ctx, u.Email).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not select from users table")
	// }
	return u, nil
}

func (s *Store) FindUser(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	st, err := s.DB.PrepareContext(ctx, "SELECT BIN_TO_UUID(id),email,passwordHash,createdAt,updatedAt FROM users WHERE email=?")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}

	if err = st.QueryRowContext(ctx, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "could not select row from users table")
	}
	return u, nil
}

func (s *Store) CheckUserExists(ctx context.Context, email string) (bool, error) {
	st, err := s.DB.PrepareContext(ctx, "SELECT EXISTS (SELECT id FROM users WHERE email=?)")
	if err != nil {
		return false, errors.Wrap(err, "could not prepare select statement")
	}
	var ok int
	if err := st.QueryRowContext(ctx, email).Scan(&ok); err != nil {
		return false, errors.Wrap(err, "could not select from users table")
	}
	return ok == 1, nil
}
