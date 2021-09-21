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
	GetUrlByID(context.Context, int64, string) (*domain.Url, error)
	GetUrls(context.Context, int, string, time.Time) ([]domain.Url, error)
	GetUrlByName(context.Context, string, string) (*domain.Url, error)
}

type Cache interface {
	Set(context.Context, string, []domain.Url, time.Duration) error
	Get(context.Context, string) ([]domain.Url, error)
	Del(context.Context, string) error
}

type Service struct {
	Store Store
	Cache Cache
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
	if ok == false {
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

func (s *Service) GetUrlByID(ctx context.Context, id int64) (*domain.Url, error) {
	ownerID := ctx.Value("userID").(string)
	url, err := s.Store.GetUrlByID(ctx, id, ownerID)
	return url, err
}
func (s *Service) GetUrlByName(ctx context.Context, name string) (*domain.Url, error) {
	ownerID := ctx.Value("userID").(string)
	url, err := s.Store.GetUrlByName(ctx, name, ownerID)
	return url, err
}

func (s *Service) GetUrls(ctx context.Context, limit int, createdAt time.Time) ([]domain.Url, error) {
	ownerID := ctx.Value("userID").(string)
	if limit < 0 {
		limit = 1
	}
	cachedUrls, err := s.Cache.Get(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if len(cachedUrls) > 0 {
		errorChan := make(chan error)
		go func() {
			urls, err := s.Store.GetUrls(ctx, limit, ownerID, cachedUrls[len(cachedUrls)-1].CreatedAt)
			if err != nil || len(urls) == 0 {
				return
			}
			err = s.Cache.Set(ctx, ownerID, urls, time.Hour)
			errorChan <- err
		}()
		err := <-errorChan
		if err != nil {
			return nil, err
		}
		return cachedUrls, nil
	}
	urls, err := s.Store.GetUrls(ctx, limit, ownerID, createdAt)
	if err != nil {
		return nil, err
	}
	go func() {
		s.Cache.Set(ctx, ownerID, urls, time.Hour)
	}()
	return urls, nil
}
