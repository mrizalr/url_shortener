package mocks

import (
	"context"

	"github.com/mrizalr/urlshortener/domain"
	"github.com/stretchr/testify/mock"
)

type UrlRepository struct {
	mock.Mock
}

func (r *UrlRepository) Create(ctx context.Context, params domain.CreateUrlParams) (int, error) {
	args := r.Mock.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (r *UrlRepository) FindByShortUrl(ctx context.Context, shortUrl string) (domain.Url, error) {
	args := r.Mock.Called(ctx, shortUrl)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (r *UrlRepository) FindByID(ctx context.Context, id int) (domain.Url, error) {
	args := r.Mock.Called(ctx, id)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (r *UrlRepository) FindAll(ctx context.Context) ([]domain.Url, error) {
	args := r.Mock.Called(ctx)
	return args.Get(0).([]domain.Url), args.Error(1)
}

func (r *UrlRepository) DeleteByID(ctx context.Context, id int) (int, error) {
	args := r.Mock.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (r *UrlRepository) IncrementURLClick(ctx context.Context, id int, count int) error {
	args := r.Mock.Called(ctx, id, count)
	return args.Error(0)
}
