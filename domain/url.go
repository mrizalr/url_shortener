package domain

import "context"

type Url struct {
	ID         int    `json:"id"`
	Url        string `json:"url"`
	ShortUrl   string `json:"short_url"`
	ClickCount int    `json:"click_count"`
	CreatedAt  int64  `json:"created_at"`
}

type CreateUrlParams struct {
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
}

type UrlRepository interface {
	Create(context.Context, CreateUrlParams) (int, error)
	FindByShortUrl(context.Context, string) (Url, error)
	FindAll(context.Context) ([]Url, error)
	DeleteByID(context.Context, int) (int, error)
}

type UrlUsecase interface {
}
