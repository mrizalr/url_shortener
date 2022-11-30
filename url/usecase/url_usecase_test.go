package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewURL(t *testing.T) {
	repoMock := new(mocks.UrlRepository)
	urlUsecase := urlUsecase{repoMock}

	urlTest := "www.github.com/mrizalr/urlshortener"
	result := domain.Url{
		ID:         1,
		Url:        fmt.Sprintf("https://%s", urlTest),
		ShortUrl:   "s5HbKw",
		ClickCount: 172,
		CreatedAt:  time.Now().Unix(),
	}

	// mock test with the case if the same random url is found in the database
	rand.Seed(time.Now().UnixNano())
	times := rand.Intn(10)
	t.Log(times)
	for i := 0; i < times; i++ {
		foundUrl := result
		var errFound error = nil

		if i == times-1 {
			foundUrl = domain.Url{}
			errFound = sql.ErrNoRows
		}

		repoMock.On("FindByShortUrl", context.Background(), mock.AnythingOfType("string")).
			Return(foundUrl, errFound).Once()
	}

	repoMock.On("Create", context.Background(), mock.AnythingOfType("domain.CreateUrlParams")).
		Return(1, nil)

	repoMock.On("FindByShortUrl", context.Background(), mock.AnythingOfType("string")).
		Return(result, nil)

	url, err := urlUsecase.CreateNewURL(context.Background(), urlTest)
	t.Log(url)

	repoMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotZero(t, url.ID)
	assert.Equal(t, result.Url, url.Url)
	assert.NotEqual(t, "", url.ShortUrl)
	assert.NotZero(t, url.CreatedAt)
}

func TestFindUrlByShort(t *testing.T) {
	repoMock := new(mocks.UrlRepository)
	urlUsecase := urlUsecase{repoMock}

	shortUrlTest := "pqS63Ns"
	result := domain.Url{
		ID:         23,
		Url:        "https://www.linkedin.com/in/mrizalr",
		ShortUrl:   shortUrlTest,
		ClickCount: 723,
		CreatedAt:  time.Date(2021, 05, 22, 18, 23, 33, 21, time.Local).Unix(),
	}

	repoMock.On("FindByShortUrl", context.Background(), shortUrlTest).Return(result, nil)

	url, err := urlUsecase.FindUrlByShort(context.Background(), shortUrlTest)
	repoMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, url.ID)
	assert.Equal(t, result.Url, url.Url)
	assert.Equal(t, result.ShortUrl, url.ShortUrl)
	assert.Equal(t, result.ClickCount, url.ClickCount)
	assert.Equal(t, result.CreatedAt, url.CreatedAt)
}

func TestFindAllUrl(t *testing.T) {
	repoMock := new(mocks.UrlRepository)
	urlUsecase := urlUsecase{repoMock}

	result := []domain.Url{
		{
			ID:         1,
			Url:        "https://www.github.com/mrizalr/urlshortener",
			ShortUrl:   "2HsEgd",
			ClickCount: 218,
			CreatedAt:  time.Now().Unix(),
		},
		{
			ID:         2,
			Url:        "https://www.linkedin.com/in/mrizalr",
			ShortUrl:   "jUHH23x",
			ClickCount: 63,
			CreatedAt:  time.Now().Unix(),
		},
	}

	repoMock.On("FindAll", context.Background()).Return(result, nil)

	urls, err := urlUsecase.FindAllUrl(context.Background())
	repoMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
}

func TestDeleteByID(t *testing.T) {
	repoMock := new(mocks.UrlRepository)
	urlUsecase := urlUsecase{repoMock}

	idTest := 1
	result := domain.Url{
		ID:         idTest,
		Url:        "https://www.linkedin.com/in/mrizalr",
		ShortUrl:   "jUHH23x",
		ClickCount: 63,
		CreatedAt:  time.Now().Unix(),
	}

	repoMock.On("DeleteByID", context.Background(), idTest).Return(idTest, nil)
	repoMock.On("FindByID", context.Background(), idTest).Return(result, nil)

	url, err := urlUsecase.DeleteByID(context.Background(), idTest)
	repoMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, url.ID)
	assert.Equal(t, result.Url, url.Url)
	assert.Equal(t, result.ShortUrl, url.ShortUrl)
	assert.Equal(t, result.ClickCount, url.ClickCount)
	assert.Equal(t, result.CreatedAt, url.CreatedAt)
}

func TestIncrementCountClick(t *testing.T) {
	repoMock := new(mocks.UrlRepository)
	urlUsecase := urlUsecase{repoMock}

	idTest := 1
	result := domain.Url{
		ID:         idTest,
		Url:        "https://www.linkedin.com/in/mrizalr",
		ShortUrl:   "jUHH23x",
		ClickCount: 63,
		CreatedAt:  time.Now().Unix(),
	}

	repoMock.On("FindByID", context.Background(), idTest).Return(result, nil)
	repoMock.On("IncrementURLClick", context.Background(), idTest, result.ClickCount+1).Return(nil)

	err := urlUsecase.IncrementClickCount(context.Background(), idTest)
	assert.NoError(t, err)
}
