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
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

func initConfig() {
	viper.SetConfigType("json")
	wd, _ := os.Getwd()
	viper.AddConfigPath(wd)
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error when read config file : " + err.Error())
	}

	log.Println("Success init config.json")
}

func init() {
	initConfig()
}

func main() {
	var (
		user     = viper.GetString("database.user")
		password = viper.GetString("database.password")
		host     = viper.GetString("database.host")
		port     = viper.GetInt("database.port")
		dbname   = viper.GetString("database.dbname")
	)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname))
	if err != nil {
		panic(err)
	}

	_mux := mux.NewRouter()

	urlRepository := repository.NewUrlRepository(db)
	urlUsecase := usecase.NewUrlUsecase(urlRepository)
	delivery.NewUrlHandler(urlUsecase, _mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("server_port")), _mux))
}
