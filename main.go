package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mrizalr/urlshortener/url/delivery"
	"github.com/mrizalr/urlshortener/url/repository"
	"github.com/mrizalr/urlshortener/url/usecase"

	_ "github.com/go-sql-driver/mysql"
)

// func initConfig() {
// 	viper.SetConfigType("json")
// 	wd, _ := os.Getwd()
// 	viper.AddConfigPath(wd)
// 	viper.SetConfigName("config")

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Fatal("Error when read config file : " + err.Error())
// 	}

// 	log.Println("Success init config.json")
// }

// func init() {
// 	initConfig()
// }

func main() {
	var (
		user     = os.Getenv("MYSQLUSER")
		password = os.Getenv("MYSQLPASSWORD")
		host     = os.Getenv("MYSQLHOST")
		port     = os.Getenv("MYSQLPORT")
		dbname   = os.Getenv("MYSQLDATABASE")
	)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname))
	if err != nil {
		panic(err)
	}

	_mux := mux.NewRouter().StrictSlash(true)
	_mux.Handle("/views/", http.FileServer(http.Dir("./views")))

	urlRepository := repository.NewUrlRepository(db)
	urlUsecase := usecase.NewUrlUsecase(urlRepository)
	delivery.NewUrlHandler(urlUsecase, _mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), _mux))
}
