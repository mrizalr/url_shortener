package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/utils"
)

type config struct {
	UrlMinLength    int
	UrlMaxLength    int
	UserIdMinLength int
	UserIdMaxLength int
}

type urlUsecase struct {
	urlRepository domain.UrlRepository
}

var _config config = config{
	UrlMinLength:    5,
	UrlMaxLength:    8,
	UserIdMinLength: 9,
	UserIdMaxLength: 11,
}

func NewUrlUsecase(urlRepository domain.UrlRepository) domain.UrlUsecase {
	return &urlUsecase{urlRepository}
}

func generateRandomURL() string {
	return utils.GetRandomURL(_config.UrlMinLength, _config.UrlMaxLength)
}

func generateRandomUserID() string {
	return utils.GetRandomURL(_config.UserIdMinLength, _config.UserIdMaxLength)
}

func (u *urlUsecase) CreateNewURL(ctx context.Context, url string) (domain.Url, error) {
	result := domain.Url{}
	if url == "" {
		return result, errors.New("validation error: url shouldn't be empty")
	}

	regex, err := regexp.Compile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	if err != nil {
		return result, err
	}

	if valid := regex.FindAllString(url, -1); len(valid) == 0 {
		return result, errors.New("validation error: url isn't valid")
	}

	if !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}

	shortUrl := generateRandomURL()
	for {
		_, err := u.urlRepository.FindByShortUrl(context.Background(), shortUrl)
		if err == sql.ErrNoRows {
			break
		}
		shortUrl = generateRandomURL()
		time.Sleep(time.Nanosecond)
	}

	var userId string
	if val := ctx.Value("user_id"); val != nil {
		userId = val.(string)
	} else {
		userId = generateRandomUserID()
	}

	params := domain.CreateUrlParams{
		Url:       url,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now().Unix(),
		UserId:    userId,
	}

	_, err = u.urlRepository.Create(context.Background(), params)
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
	url, err := u.urlRepository.FindByID(context.Background(), id)
	if err != nil {
		return url, err
	}

	_, err = u.urlRepository.DeleteByID(ctx, id)
	return url, err
}

func (u *urlUsecase) IncrementClickCount(ctx context.Context, id int) error {
	url, err := u.urlRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.urlRepository.IncrementURLClick(ctx, id, url.ClickCount+1)
	if err != nil {
		return err
	}

	return nil
}

func (u *urlUsecase) GetLastUrlCreated(ctx context.Context, userId string) ([]domain.Url, error) {
	urls, err := u.urlRepository.GetLastUrlCreated(ctx, userId)
	return urls, err
}
