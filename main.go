package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shvz0/mmscrap/db"
	"github.com/shvz0/mmscrap/mmscrappers"
	"github.com/shvz0/mmscrap/server"
)

func main() {
	godotenv.Load(".env")

	conn := db.Connect()

	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0777)
	}

	file, err := os.OpenFile("./logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

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
