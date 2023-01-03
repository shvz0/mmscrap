package main

import (
	"flag"

	"github.com/shvz0/mmscrap/db"
	"github.com/shvz0/mmscrap/mmscrappers"
	"github.com/shvz0/mmscrap/server"
)

func main() {

	// n24 := mmscrappers.NewNews24()

	// n24.ArticleListToday()

	migrate := flag.Bool("migrate", false, "Perform DB migration")
	serve := flag.Bool("serve", false, "Run server")
	parse := flag.Bool("parse", false, "Run server")

	flag.Parse()

	if *migrate {
		db.Migrate()
	}

	if *parse {
		mmscrappers.SaveTodaysArticlesToDB(db.Db)
	}

	if *serve {
		server.Serve()
	}
}
