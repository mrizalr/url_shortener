package mocks

import (
	"context"

	"github.com/mrizalr/urlshortener/domain"
	"github.com/stretchr/testify/mock"
)

type UrlUsecase struct {
	mock.Mock
}

func (u *UrlUsecase) CreateNewURL(ctx context.Context, url string) (domain.Url, error) {
	args := u.Mock.Called(ctx, url)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (u *UrlUsecase) FindUrlByShort(ctx context.Context, shortUrl string) (domain.Url, error) {
	args := u.Mock.Called(ctx, shortUrl)
	return args.Get(0).(domain.Url), args.Error(1)
}

func (u *UrlUsecase) FindAllUrl(ctx context.Context) ([]domain.Url, error) {
	args := u.Mock.Called(ctx)
	return args.Get(0).([]domain.Url), args.Error(1)
}

func (u *UrlUsecase) DeleteByID(ctx context.Context, ID int) (domain.Url, error) {
	args := u.Mock.Called(ctx, ID)
	return args.Get(0).(domain.Url), args.Error(1)
}
