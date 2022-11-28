package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/utils"
)

type config struct {
	UrlMinLength int
	UrlMaxLength int
}

type urlUsecase struct {
	urlRepository domain.UrlRepository
}

var _config config = config{
	UrlMinLength: 5,
	UrlMaxLength: 8,
}

func NewUrlUsecase(urlRepository domain.UrlRepository) domain.UrlUsecase {
	return &urlUsecase{urlRepository}
}

func generateRandom() string {
	return utils.GetRandomURL(_config.UrlMinLength, _config.UrlMaxLength)
}

func (u *urlUsecase) CreateNewURL(ctx context.Context, url string) (domain.Url, error) {
	result := domain.Url{}
	shortUrl := generateRandom()
	for {
		_, err := u.urlRepository.FindByShortUrl(context.Background(), shortUrl)
		if err == sql.ErrNoRows {
			break
		}
		shortUrl = generateRandom()
		time.Sleep(time.Nanosecond)
	}

	params := domain.CreateUrlParams{
		Url:      url,
		ShortUrl: shortUrl,
	}

	_, err := u.urlRepository.Create(context.Background(), params)
	if err != nil {
		return result, err
	}

	result, err = u.urlRepository.FindByShortUrl(context.Background(), shortUrl)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *urlUsecase) FindUrlByShort(ctx context.Context, shortUrl string) (domain.Url, error) {
	url, err := u.urlRepository.FindByShortUrl(ctx, shortUrl)
	return url, err
}

func (u *urlUsecase) FindAllUrl(ctx context.Context) ([]domain.Url, error) {
	urls, err := u.urlRepository.FindAll(ctx)
	return urls, err
}

func (u *urlUsecase) DeleteByID(ctx context.Context, id int) (domain.Url, error) {
	url := domain.Url{}
	id, err := u.urlRepository.DeleteByID(ctx, id)
	if err != nil {
		return url, err
	}

	url, err = u.urlRepository.FindByID(context.Background(), id)
	return url, err
}
