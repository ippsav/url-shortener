package mysql

import (
	"context"
	"time"
	"url-shortner/domain"

	"github.com/pkg/errors"
)

func (s *Store) CreateUrlsTable(ctx context.Context) error {
	table := `CREATE TABLE IF NOT EXISTS urls(
   id INT PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR(16) ,
   redirectTo  VARCHAR(25) ,
   ownerID BINARY(16),
   createdAt DATETIME,
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

func (s *Store) GetUrlByID(ctx context.Context, id int64, ownerID string) (*domain.Url, error) {
	st, err := s.DB.PrepareContext(ctx, "SELECT id,name,redirectTo,BIN_TO_UUID(ownerID),createdAt FROM urls WHERE id=? AND ownerID=UUID_TO_BIN(?)")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	url := &domain.Url{}
	err = st.QueryRowContext(ctx, id, ownerID).Scan(&url.ID, &url.Name, &url.RedirectTo, &url.OwnerID, &url.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not select from urls")
	}
	return url, nil
}
func (s *Store) GetUrlByName(ctx context.Context, name, ownerID string) (*domain.Url, error) {
	st, err := s.DB.PrepareContext(ctx, "SELECT id,name,redirectTo,BIN_TO_UUID(ownerID),createdAt FROM urls WHERE name=? AND ownerID=UUID_TO_BIN(?)")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	u := &domain.Url{}
	err = st.QueryRowContext(ctx, name, ownerID).Scan(&u.ID, &u.Name, &u.RedirectTo, &u.OwnerID, &u.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not select from urls")
	}
	return u, nil
}

func (s *Store) GetUrls(ctx context.Context, limit int, ownerID string, createdAt time.Time) ([]domain.Url, error) {
	//AND createdAt < ? ORDER BY createdAt DESC LIMIT ?
	st, err := s.DB.PrepareContext(ctx, "SELECT id,name,redirectTo,BIN_TO_UUID(ownerID),createdAt FROM urls WHERE ownerID=UUID_TO_BIN(?) AND createdAt < ? ORDER BY createdAt DESC LIMIT ? ")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	rows, err := st.QueryContext(ctx, ownerID, createdAt, limit)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "could not select from urls")
	}
	urls := make([]domain.Url, 0)
	for rows.Next() {
		u := &domain.Url{}
		err := rows.Scan(&u.ID, &u.Name, &u.RedirectTo, &u.OwnerID, &u.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "could not read the current row")
		}
		urls = append(urls, *u)
	}
	return urls, nil
}
