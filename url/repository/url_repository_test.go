package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrizalr/urlshortener/db/queries"
	"github.com/mrizalr/urlshortener/domain"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(err)
	}
	return db, mock
}

func TestCreateURL(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	params := domain.CreateUrlParams{
		Url:       "www.github.com/mrizalr/urlshortener",
		ShortUrl:  "xhYsg23",
		CreatedAt: time.Now().Unix(),
		UserId:    "hd6GSj2bGFSh112",
	}

	mock.ExpectExec(queries.InsertURL).
		WithArgs(params.Url, params.ShortUrl, params.CreatedAt, params.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := urlRepository{db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := repo.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestFindByShortUrl(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	params := domain.Url{
		ID:         1,
		Url:        "www.github.com/mrizalr/urlshortener",
		ShortUrl:   "xh52VsC",
		ClickCount: 162,
		CreatedAt:  time.Now().Unix(),
		UserId:     "hd6GSj2bGFSh112",
	}

	rows := mock.NewRows([]string{"id", "url", "short_url", "click_count", "created_at", "user_id"}).
		AddRow(params.ID, params.Url, params.ShortUrl, params.ClickCount, params.CreatedAt, params.UserId)
	mock.ExpectQuery(queries.FindByShort).WithArgs(params.ShortUrl).WillReturnRows(rows)

	repo := urlRepository{db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	url, err := repo.FindByShortUrl(ctx, params.ShortUrl)

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, params.ID, url.ID)
	assert.Equal(t, params.Url, url.Url)
	assert.Equal(t, params.ShortUrl, url.ShortUrl)
	assert.Equal(t, params.ClickCount, url.ClickCount)
	assert.Equal(t, params.CreatedAt, url.CreatedAt)
	assert.Equal(t, params.UserId, url.UserId)
}

func TestFindAll(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	params := []domain.Url{
		{
			ID:         1,
			Url:        "www.github.com/mrizalr/urlshortener",
			ShortUrl:   "2HsEgd",
			ClickCount: 218,
			CreatedAt:  time.Now().Unix(),
			UserId:     "asHJS661GGsdGa",
		},
		{
			ID:         2,
			Url:        "www.linkedin.com/in/mrizalr",
			ShortUrl:   "jUHH23x",
			ClickCount: 63,
			CreatedAt:  time.Now().Unix(),
			UserId:     "asHJS661GGsdGa",
		},
	}

	rows := mock.NewRows([]string{"id", "url", "short_url", "click_count", "created_at", "user_id"})
	for _, param := range params {
		rows.AddRow(param.ID, param.Url, param.ShortUrl, param.ClickCount, param.CreatedAt, param.UserId)
	}

	mock.ExpectQuery(queries.FindAll).WillReturnRows(rows)

	repo := urlRepository{db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	urls, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, urls)
	assert.Len(t, urls, 2)
}

func TestDeleteByID(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	params := domain.Url{
		ID:         1,
		Url:        "www.github.com/mrizalr/urlshortener",
		ShortUrl:   "ofJA32",
		ClickCount: 749,
		CreatedAt:  time.Now().Unix(),
		UserId:     "asHJS661GGsdGa",
	}

	mock.ExpectExec(queries.DeleteByID).WithArgs(params.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := urlRepository{db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := repo.DeleteByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestIncrementUrl(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	params := domain.Url{
		ID:         1,
		Url:        "www.github.com/mrizalr/urlshortener",
		ShortUrl:   "ofJA32",
		ClickCount: 12,
		CreatedAt:  time.Now().Unix(),
		UserId:     "asHJS661GGsdGa",
	}

	mock.ExpectExec(queries.IncrementClickCount).
		WithArgs(params.ClickCount+1, params.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := urlRepository{db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := repo.IncrementURLClick(ctx, 1, 13)
	assert.NoError(t, err)
}

func TestGetLastUrlCreated(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	result := []domain.Url{{
		ID:         1,
		Url:        "www.github.com/1",
		ShortUrl:   "j723HS",
		ClickCount: 192,
		CreatedAt:  time.Now().Add(-2 * time.Hour).Unix(),
		UserId:     "sf7y612aHAsd",
	}, {
		ID:         2,
		Url:        "www.github.com/3",
		ShortUrl:   "j723HS",
		ClickCount: 192,
		CreatedAt:  time.Now().Add(-1 * time.Hour).Unix(),
		UserId:     "sf7y612aHAsd",
	}, {
		ID:         3,
		Url:        "www.github.com",
		ShortUrl:   "j723HS",
		ClickCount: 192,
		CreatedAt:  time.Now().Unix(),
		UserId:     "sf7y612aHAsd",
	}}

	rows := mock.NewRows([]string{"id", "url", "short_url", "clicked_count", "created_at", "user_id"})
	for _, val := range result {
		rows.AddRow(val.ID, val.Url, val.ShortUrl, val.ClickCount, val.CreatedAt, val.UserId)
	}

	mock.ExpectQuery(queries.GetLastUrl).WithArgs(result[0].UserId).WillReturnRows(rows)

	repo := urlRepository{db}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	urls, err := repo.GetLastUrlCreated(ctx, result[0].UserId)
	assert.NoError(t, err)

	assert.Len(t, urls, 3)
}
