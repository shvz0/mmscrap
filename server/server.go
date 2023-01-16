package server

import (
	"log"
	"net/http"

	"github.com/shvz0/mmscrap/db"

	"gorm.io/gorm"
)

var Db *gorm.DB

func initServer() {
	Db = db.Connect()
}

func Serve() {
	initServer()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./public"))

	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	mux.Handle("/home", MainPageHandler{})
	mux.Handle("/", MainPageHandler{})

	mux.Handle("/delta", DeltaHandler{})
	mux.Handle("/mendenhall", MendenhallHandler{})
	mux.Handle("/chisquared", ChiSquredHandler{})

	log.Print("Listening...")
	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		panic(err)
	}
}
