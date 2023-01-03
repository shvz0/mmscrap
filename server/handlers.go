package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mmscrap/mmscrappers"
	"net/http"
	"os"
	"strings"
)

type MainPageHandler struct {
}

func (h MainPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, err := os.OpenFile("./templates/index.html", os.O_RDONLY, 0777)

	var all []mmscrappers.Article
	var buf strings.Builder
	_, err = io.Copy(&buf, f)

	alldb := Db.Order("date desc").Find(&all)
	rows, err := alldb.Rows()

	if err != nil {
		fmt.Println(err)
	}

	rows.Close()

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(string(buf.String()))
	check(err)

	data := struct {
		Title      string
		DateFormat string
		Items      []mmscrappers.Article
	}{
		Title:      "Home",
		DateFormat: "01.02.2006 15:04",
		Items:      all,
	}

	err = t.Execute(w, data)
	check(err)
}
