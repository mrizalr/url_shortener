package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/mrizalr/urlshortener/domain"
	"github.com/mrizalr/urlshortener/utils"
)

type UrlHandler struct {
	urlUsecase domain.UrlUsecase
}

func NewUrlHandler(urlUsecase domain.UrlUsecase, m *mux.Router) {
	handler := UrlHandler{urlUsecase}
	m.HandleFunc("/", handler.HomeHandler).Methods("GET")
	m.HandleFunc("/{short}", handler.getUrlByShort)

	router_v1 := m.PathPrefix("/api/v1/url").Subrouter()

	router_v1.Path("/").HandlerFunc(handler.getAllUrl).Methods("GET")
	router_v1.Path("/create").HandlerFunc(handler.createNewUrlShortener).Methods("POST")
	router_v1.Path("/{id}").HandlerFunc(handler.deleteUrlByID).Methods("DELETE")

}

func (h *UrlHandler) HomeHandler(res http.ResponseWriter, req *http.Request) {
	filePath := path.Join("views", "index.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		res.Write([]byte("Bad gateway"))
		return
	}

	errs := ""
	cardTemplate := `<div class="card">
						<div class="short-link">%s</div>
						<div class="web-title">%s</div>
						<div class="web-url">%s</div>
						<div class="created-date">%s</div>
						<div class="clicked-count">
							<span class="count">%d</span>
							clicks
						</div>
					</div>`

	cards := ""
	if storedCookie, _ := req.Cookie("user_id"); storedCookie != nil {
		userId := storedCookie.Value
		urls, err := h.urlUsecase.GetLastUrlCreated(context.Background(), userId)
		if err != nil {
			errs += err.Error() + "|"
		}

		for _, url := range urls {
			var webtitle string
			res, err := http.Get(url.Url)
			if err != nil {
				errs += err.Error() + "|"
			}
			defer res.Body.Close()

			if res.StatusCode == 200 {
				doc, _ := goquery.NewDocumentFromReader(res.Body)
				webtitle = doc.Find("title").Text()
			} else {
				if err != nil {
					errs += fmt.Sprintf("status code %d", res.StatusCode) + "|"
				}
			}

			createdAt := time.Unix(url.CreatedAt, 0)
			cards += fmt.Sprintf(cardTemplate, url.ShortUrl, webtitle, url.Url, createdAt, url.ClickCount)
		}
		res.Write([]byte(userId))
	}

	data := map[string]interface{}{
		"cards": cards,
	}

	err = tmpl.Execute(res, data)
	res.Write([]byte(cards))
	if err != nil {
		errs += err.Error() + "|"
		res.Write([]byte(errs))
	}
}

func (h *UrlHandler) createNewUrlShortener(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")

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

	ctx := context.Background()
	if storedCookie, _ := req.Cookie("user_id"); storedCookie != nil {
		userId := storedCookie.Value
		ctx = context.WithValue(ctx, "user_id", userId)
	}

	url, err := h.urlUsecase.CreateNewURL(ctx, requestBody.Url)
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

	http.SetCookie(res, &http.Cookie{
		Name:  "user_id",
		Value: url.UserId,
	})

	utils.FormatResponse(res, &utils.ResponseSuccessParams{
		Code:   http.StatusCreated,
		Status: "Success Created",
		Data:   url,
	})
}

func (h *UrlHandler) getAllUrl(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")

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
	res.Header().Set("Access-Control-Allow-Origin", "*")

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
		filePath := path.Join("views", "404.html")
		tmpl, err := template.ParseFiles(filePath)
		if err != nil {
			res.Write([]byte("Bad gateway"))
		}

		err = tmpl.Execute(res, nil)
		if err != nil {
			res.Write([]byte("Bad gateway"))
		}

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
