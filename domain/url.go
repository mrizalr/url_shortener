package domain

import "context"

type Url struct {
	ID         int    `json:"id"`
	Url        string `json:"url"`
	ShortUrl   string `json:"short_url"`
	ClickCount int    `json:"click_count"`
	CreatedAt  int64  `json:"created_at"`
	UserId     string `json:"user_id"`
}

type CreateUrlParams struct {
	Url       string `json:"url"`
	ShortUrl  string `json:"short_url"`
	CreatedAt int64  `json:"created_at"`
	UserId    string `json:"user_id"`
}

type UrlRepository interface {
	Create(context.Context, CreateUrlParams) (int, error)
	FindByShortUrl(context.Context, string) (Url, error)
	FindByID(context.Context, int) (Url, error)
	FindAll(context.Context) ([]Url, error)
	DeleteByID(context.Context, int) (int, error)
	IncrementURLClick(context.Context, int, int) error
}

type UrlUsecase interface {
	CreateNewURL(context.Context, string) (Url, error)
	FindUrlByShort(context.Context, string) (Url, error)
	FindAllUrl(context.Context) ([]Url, error)
	DeleteByID(context.Context, int) (Url, error)
	IncrementClickCount(context.Context, int) error
}
