package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/shvz0/mmscrap/db"
	"github.com/shvz0/mmscrap/mmscrappers"
	"github.com/shvz0/mmscrap/server"
)

func main() {
	godotenv.Load(".env")

	conn := db.Connect()

	migrate := flag.Bool("migrate", false, "Perform DB migration")
	serve := flag.Bool("serve", false, "Run server")
	parse := flag.Bool("parse", false, "Run parsing")

	flag.Parse()

	if *migrate {
		db.Migrate()
	}

	if *parse {
		mmscrappers.SaveTodaysArticlesToDB(conn)
	}

	if *serve {
		server.Serve()
	}
}
