package repository

import (
	"context"
	"database/sql"

	"github.com/mrizalr/urlshortener/db/queries"
	"github.com/mrizalr/urlshortener/domain"
)

type urlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) domain.UrlRepository {
	return &urlRepository{db}
}

// Inserting new shortener url data to urls table
// Receiving context, and CreateURLParams as parameter
// Returning inserted url_id (int) if success, and error if failed

func (r *urlRepository) Create(ctx context.Context, params domain.CreateUrlParams) (int, error) {
	sqlRes, err := r.db.ExecContext(ctx, queries.InsertURL, params.Url, params.ShortUrl)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}

// Fetch one url data from urls table
// Receiving context, and shortUrl (string) as parameter
// Returning url data (domain.Url) if success, and error if failed

func (r *urlRepository) FindByShortUrl(ctx context.Context, shortUrl string) (domain.Url, error) {
	url := domain.Url{}
	err := r.db.QueryRowContext(ctx, queries.FindByShort, shortUrl).
		Scan(&url.ID, &url.Url, &url.ShortUrl, &url.ClickCount, &url.CreatedAt)

	if err != nil {
		return url, err
	}

	return url, nil
}

// Fetch all url data from urls table
// Receiving context as parameter
// Returning url data ([] domain.Url) if success, and error if failed

func (r *urlRepository) FindAll(ctx context.Context) ([]domain.Url, error) {
	urls := []domain.Url{}
	rows, err := r.db.QueryContext(ctx, queries.FindAll)
	if err != nil {
		return urls, err
	}
	defer rows.Close()

	for rows.Next() {
		url := domain.Url{}
		err = rows.Scan(&url.ID, &url.Url, &url.ShortUrl, &url.ClickCount, &url.CreatedAt)
		if err != nil {
			return urls, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}

// Delete one url data from urls table
// Receiving context, and url_id (int) as parameter
// Returning deleted url_id (int) if success, and error if failed

func (r *urlRepository) DeleteByID(ctx context.Context, ID int) (int, error) {
	sqlRes, err := r.db.ExecContext(ctx, queries.DeleteByID, ID)
	if err != nil {
		return 0, err
	}

	lastDeletedId, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastDeletedId), nil
}
