package delivery

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewUrlHandler(t *testing.T) {
	mockUsecase := new(mocks.UrlUsecase)
	usecaseResult := domain.Url{
		ID:         1,
		Url:        "www.github.com/mrizalr",
		ShortUrl:   "h52GbxA",
		ClickCount: 0,
		CreatedAt:  time.Now().Unix(),
	}

	mockUsecase.On("CreateNewURL", context.Background(), mock.AnythingOfType("string")).Return(usecaseResult, nil)

	reqJson := fmt.Sprintf(`{"url":"%s"}`, usecaseResult.Url)
	reqBody := bytes.NewReader([]byte(reqJson))

	req := httptest.NewRequest("POST", "/api/v1/url/create", reqBody)
	res := httptest.NewRecorder()

	handler := UrlHandler{mockUsecase}
	handler.createNewUrlShortener(res, req)

	result := res.Result()
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	assert.NoError(t, err)

	expect := fmt.Sprintf(`
	{
		"status_code":201,
		"status":"Success Created",
		"data":{
			"id":1,
			"url":"www.github.com/mrizalr",
			"short_url":"h52GbxA",
			"click_count":0,
			"created_at":%d
		}
	}`, time.Now().Unix())

	mockUsecase.AssertExpectations(t)
	assert.JSONEq(t, expect, string(resultBody))
}

func TestGetAllUrlHandler(t *testing.T) {
	mockUsecase := new(mocks.UrlUsecase)
	usecaseResult := []domain.Url{
		{
			ID:         1,
			Url:        "www.github.com/mrizalr",
			ShortUrl:   "h52GbxA",
			ClickCount: 163,
			CreatedAt:  time.Date(2021, 03, 23, 12, 13, 32, 43, time.Local).Unix(),
		}, {
			ID:         2,
			Url:        "www.linkedin.com/in/mrizalr",
			ShortUrl:   "hJS62h",
			ClickCount: 123,
			CreatedAt:  time.Date(2021, 03, 23, 12, 13, 32, 43, time.Local).Unix(),
		},
	}

	mockUsecase.On("FindAllUrl", context.Background()).Return(usecaseResult, nil)

	req := httptest.NewRequest("GET", "/api/v1/url/", nil)
	res := httptest.NewRecorder()

	handler := UrlHandler{mockUsecase}
	handler.getAllUrl(res, req)

	result := res.Result()
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	assert.NoError(t, err)

	expect := fmt.Sprintf(`
	{
		"status_code":200,
		"status":"Success",
		"data":[{
			"id":1,
			"url":"www.github.com/mrizalr",
			"short_url":"h52GbxA",
			"click_count":163,
			"created_at":%d
		},
		{
			"id":2,
			"url":"www.linkedin.com/in/mrizalr",
			"short_url":"hJS62h",
			"click_count":123,
			"created_at":%d
		}]
	}`, usecaseResult[0].CreatedAt, usecaseResult[1].CreatedAt)

	mockUsecase.AssertExpectations(t)
	assert.JSONEq(t, expect, string(resultBody))
}

func TestDeleteUrlByIDHandler(t *testing.T) {
	mockUsecase := new(mocks.UrlUsecase)
	usecaseResult := domain.Url{
		ID:         1,
		Url:        "www.github.com/mrizalr",
		ShortUrl:   "h52GbxA",
		ClickCount: 163,
		CreatedAt:  time.Date(2021, 03, 23, 12, 13, 32, 43, time.Local).Unix(),
	}

	mockUsecase.On("DeleteByID", context.Background(), usecaseResult.ID).Return(usecaseResult, nil)

	req := httptest.NewRequest("DELETE", "/api/v1/url/1", nil)
	res := httptest.NewRecorder()

	params := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, params)

	handler := UrlHandler{mockUsecase}
	handler.deleteUrlByID(res, req)

	result := res.Result()
	defer result.Body.Close()

	resultBody, err := io.ReadAll(result.Body)
	assert.NoError(t, err)

	expect := fmt.Sprintf(`
	{
		"status_code":200,
		"status":"Success",
		"data":{
			"id":1,
			"url":"www.github.com/mrizalr",
			"short_url":"h52GbxA",
			"click_count":163,
			"created_at":%d
		}
	}`, usecaseResult.CreatedAt)

	mockUsecase.AssertExpectations(t)
	assert.JSONEq(t, expect, string(resultBody))
}

func TestGetUrl(t *testing.T) {
	mockUsecase := new(mocks.UrlUsecase)
	usecaseResult := domain.Url{
		ID:         1,
		Url:        "https://www.google.com",
		ShortUrl:   "ha51Fad",
		ClickCount: 23,
		CreatedAt:  time.Now().Unix(),
	}

	mockUsecase.On("FindUrlByShort", context.Background(), mock.AnythingOfType("string")).
		Return(usecaseResult, nil)
	mockUsecase.On("IncrementClickCount", context.Background(), usecaseResult.ID).Return(nil)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/url/%s", usecaseResult.ShortUrl), nil)
	res := httptest.NewRecorder()

	params := map[string]string{"short": usecaseResult.ShortUrl}
	req = mux.SetURLVars(req, params)

	handler := UrlHandler{mockUsecase}
	handler.getUrlByShort(res, req)

	result := res.Result()
	redirectUrl, err := result.Location()

	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
	assert.Equal(t, usecaseResult.Url, fmt.Sprintf("%s://%s", redirectUrl.Scheme, redirectUrl.Host))
}
