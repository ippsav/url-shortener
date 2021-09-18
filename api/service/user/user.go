package user

import (
	"context"
	"errors"
	"url-shortner/domain"
)

type Store interface {
	CreateUser(context.Context, *domain.User) (*domain.User, error)
	FindUser(context.Context, string) (*domain.User, error)
}

type Service struct {
	Store Store
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	u := &domain.User{}
	u.Email = email
	if err := u.HashPassword(password); err != nil {
		return nil, err
	}
	u, err := s.Store.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
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
