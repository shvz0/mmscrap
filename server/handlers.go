package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/shvz0/mmscrap/mmscrappers"
)

type MainPageHandler struct {
}

func (h MainPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var all []mmscrappers.Article

	alldb := Db.Order("date desc").Find(&all)
	rows, err := alldb.Rows()

	if err != nil {
		log.Printf("DB error: %v", err)
		responseServerError(w, err)
		return
	}

	rows.Close()

	if err != nil {
		log.Println(err)
		responseServerError(w, err)
		return
	}

	data := struct {
		Title      string
		DateFormat string
		Items      []mmscrappers.Article
	}{
		Title:      "Home",
		DateFormat: "01.02.2006 15:04",
		Items:      all,
	}

	t, err := getTemplate("index.html")

	if err != nil {
		log.Println(err)
		responseServerError(w, err)
		return
	}

	err = t.Execute(w, data)

	if err != nil {
		log.Println(err)
		responseServerError(w, err)
		return
	}
}

func responseServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	t, _ := getTemplate("errors/500.html")
	t.Execute(w, err)
}

func getTemplate(templatePath string) (*template.Template, error) {
	var buf strings.Builder

	f, err := os.OpenFile("./templates/"+templatePath, os.O_RDONLY, 0777)

	if err != nil {
		log.Printf("IO error: %v", err)
		return nil, err
	}

	_, err = io.Copy(&buf, f)

	if err != nil {
		log.Printf("IO error: %v", err)
		return nil, err
	}

	t, err := template.New("webpage").Parse(string(buf.String()))

	if err != nil {
		log.Printf("IO error: %v", err)
		return nil, err
	}

	return t, nil
}
