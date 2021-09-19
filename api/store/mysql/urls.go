package mysql

import (
	"context"
	"url-shortner/domain"

	"github.com/pkg/errors"
)

func (s *Store) CreateUrlsTable(ctx context.Context) error {
	table := `CREATE TABLE IF NOT EXISTS urls(
   id INT PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR(16) UNIQUE,
   redirectTo  VARCHAR(25) UNIQUE,
   ownerID BINARY(16),
   createdAt DATETIME,
   UNIQUE(name,redirectTo),
   FOREIGN KEY(ownerID) REFERENCES users(id)
  )`
	_, err := s.DB.ExecContext(ctx, table)
	return err
}

func (s *Store) CreateUrl(ctx context.Context, url *domain.Url) (*domain.Url, error) {
	st, err := s.DB.PrepareContext(ctx, "INSERT INTO urls(name,redirectTo,ownerID,createdAt) values(?,?,UUID_TO_BIN(?),?)")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	res, err := st.ExecContext(ctx, url.Name, url.RedirectTo, url.OwnerID, url.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert into urls")
	}
	url.ID, err = res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "could not get the last inserted id")
	}
	return url, nil
}

func (s *Store) CheckUrlExists(ctx context.Context, name, redirectTo, ownerID string) (bool, error) {
	st, err := s.DB.PrepareContext(ctx, "SELECT EXISTS(SELECT id FROM urls WHERE name=? OR redirectTo=? AND ownerID=UUID_TO_BIN(?))")
	if err != nil {
		return false, errors.Wrap(err, "could not prepare select statement")
	}
	var ok int
	err = st.QueryRowContext(ctx, name, redirectTo, ownerID).Scan(&ok)
	if err != nil {
		return false, errors.Wrap(err, "could not select from urls")
	}
	return ok == 1, nil
}

func (s *Store) GetUrl(ctx context.Context, name, ownerID string) (*domain.Url, error) {
	st, err := s.DB.PrepareContext(ctx, "SELECT id,name,redirectTo,BIN_TO_UUID(ownerID),createdAt FROM urls WHERE name=? AND ownerID=UUID_TO_BIN(?)")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	url := &domain.Url{}
	err = st.QueryRowContext(ctx, name, ownerID).Scan(&url.ID, &url.Name, &url.RedirectTo, &url.OwnerID, &url.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not select from urls")
	}
	return url, nil
}
