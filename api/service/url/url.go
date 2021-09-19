package url

import (
	"context"
	"errors"
	"time"
	"url-shortner/domain"
)

type Store interface {
	CreateUrl(context.Context, *domain.Url) (*domain.Url, error)
	CheckUrlExists(context.Context, string, string, string) (bool, error)
	GetUrl(context.Context, string, string) (*domain.Url, error)
}

type Service struct {
	Store Store
}

func (s *Service) CreateUrl(ctx context.Context, name, redirectTo string) (*domain.Url, error) {
	ownerID := ctx.Value("userID").(string)
	url := &domain.Url{
		Name:       name,
		RedirectTo: redirectTo,
		OwnerID:    ownerID,
		CreatedAt:  time.Now(),
	}
	ok := url.Validate()
	if !ok {
		return nil, errors.New("bad url format")
	}
	url, err := s.Store.CreateUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (s *Service) CheckUrlExists(ctx context.Context, name, redirectTo string) (bool, error) {
	ownerID := ctx.Value("userID").(string)
	ok, err := s.Store.CheckUrlExists(ctx, name, redirectTo, ownerID)
	return ok, err
}
func (s *Service) GetUrl(ctx context.Context, name string) (*domain.Url, error) {
	ownerID := ctx.Value("userID").(string)
	url, err := s.Store.GetUrl(ctx, name, ownerID)
	return url, err
}
