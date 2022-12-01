package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/utils"
)

type UrlHandler struct {
	urlUsecase domain.UrlUsecase
}

func NewUrlHandler(urlUsecase domain.UrlUsecase, m *mux.Router) {
	handler := UrlHandler{urlUsecase}
	router_v1 := m.PathPrefix("/api/v1/url").Subrouter()

	router_v1.Path("/").HandlerFunc(handler.getAllUrl).Methods("GET")
	router_v1.Path("/create").HandlerFunc(handler.createNewUrlShortener).Methods("POST")
	router_v1.Path("/{id}").HandlerFunc(handler.deleteUrlByID).Methods("DELETE")
	router_v1.Path("/{short}").HandlerFunc(handler.getUrlByShort).Methods("GET")
}

func (h *UrlHandler) createNewUrlShortener(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad gateway",
			Errors: []string{err.Error()},
		})
		return
	}
	defer req.Body.Close()

	requestBody := struct {
		Url string `json:"url"`
	}{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadRequest,
			Status: "Bad request",
			Errors: []string{
				"error while parsing json",
			},
		})
		return
	}

	url, err := h.urlUsecase.CreateNewURL(context.Background(), requestBody.Url)
	if err != nil {
		errorParams := utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad Gateway",
			Errors: []string{err.Error()},
		}

		if strings.Contains(strings.ToLower(err.Error()), "validation") {
			errorParams.Code = http.StatusBadRequest
			errorParams.Status = "Bad request"
			errorParams.Errors = []string{err.Error()}
		}

		utils.FormatResponse(res, &errorParams)
		return
	}

	utils.FormatResponse(res, &utils.ResponseSuccessParams{
		Code:   http.StatusCreated,
		Status: "Success Created",
		Data:   url,
	})
}

func (h *UrlHandler) getAllUrl(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	urls, err := h.urlUsecase.FindAllUrl(context.Background())
	if err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad gateway",
			Errors: []string{err.Error()},
		})
		return
	}

	utils.FormatResponse(res, &utils.ResponseSuccessParams{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   urls,
	})
}

func (h *UrlHandler) deleteUrlByID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	urlIdStr := mux.Vars(req)["id"]
	urlId, err := strconv.Atoi(urlIdStr)
	if err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadRequest,
			Status: "Bad request",
			Errors: []string{"url id isn't valid"},
		})
		return
	}

	url, err := h.urlUsecase.DeleteByID(context.Background(), urlId)
	if err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad gateway",
			Errors: []string{err.Error()},
		})
		return
	}

	utils.FormatResponse(res, &utils.ResponseSuccessParams{
		Code:   http.StatusOK,
		Status: "Success",
		Data:   url,
	})
}

func (h *UrlHandler) getUrlByShort(res http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	xforward := req.Header.Get("X-Forwarded-For")
	fmt.Println("IP : ", ip)
	fmt.Println("X-Forwarded-For : ", xforward)

	shortUrl := mux.Vars(req)["short"]

	urlChan := make(chan domain.Url, 5)
	errChan := make(chan error, 5)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		url, err := h.urlUsecase.FindUrlByShort(context.Background(), shortUrl)
		urlChan <- url
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad gateway",
			Errors: []string{err.Error()},
		})
		return
	}

	url := <-urlChan

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.urlUsecase.IncrementClickCount(context.Background(), url.ID)
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		utils.FormatResponse(res, &utils.ResponseErrorParams{
			Code:   http.StatusBadGateway,
			Status: "Bad gateway",
			Errors: []string{err.Error()},
		})
		return
	}

	wg.Wait()
	http.Redirect(res, req, url.Url, http.StatusTemporaryRedirect)
}
