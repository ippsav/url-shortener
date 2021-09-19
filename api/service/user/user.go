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
	CheckUserExistsWithID(context.Context, string) (bool, error)
}

type Service struct {
	Store Store
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	u := &domain.User{
		ID:        uuid.NewString(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ok := u.Validate()
	if !ok {
		return nil, errors.New("wrong email format")
	}
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
	return ok, err
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
