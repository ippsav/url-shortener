package user

import (
	"context"
	"errors"
	"time"
	"url-shortner/domain"

	"github.com/google/uuid"
)

type Store interface {
	CreateUser(context.Context, *domain.User) (*domain.User, error)
	FindUser(context.Context, string) (*domain.User, error)
	CheckUserExists(context.Context, string) (bool, error)
}

type Service struct {
	Store Store
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	u := &domain.User{}
	u.Email = email
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if err := u.HashPassword(password); err != nil {
		return nil, err
	}
	u, err := s.Store.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
func (s *Service) CheckUserExists(ctx context.Context, email string) (bool, error) {
	ok, err := s.Store.CheckUserExists(ctx, email)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *Service) FindUser(ctx context.Context, email, password string) (*domain.User, error) {
	u, err := s.Store.FindUser(ctx, email)
	if err != nil {
		return nil, err
	}
	if ok := u.CheckPasswordHash(password); ok != true {
		return nil, errors.New("wrong password")
	}
	return u, nil
}
