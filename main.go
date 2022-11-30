package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrizalr/urlshortener/url/delivery"
	"github.com/mrizalr/urlshortener/url/repository"
	"github.com/mrizalr/urlshortener/url/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_mux := mux.NewRouter()
	db, err := sql.Open("mysql", "root:secret@tcp(localhost:2252)/url_short")
	if err != nil {
		panic(err)
	}

	urlRepository := repository.NewUrlRepository(db)
	urlUsecase := usecase.NewUrlUsecase(urlRepository)
	delivery.NewUrlHandler(urlUsecase, _mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), _mux))
}
