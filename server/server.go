package server

import (
	"log"
	"net/http"

	"github.com/shvz0/mmscrap/db"

	"gorm.io/gorm"
)

var dsn string
var Db *gorm.DB

func initServer() {
	Db = db.Db
}

func Serve() {
	initServer()
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)
	mux.Handle("/home", MainPageHandler{})

	log.Print("Listening...")
	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		panic(err)
	}
}
