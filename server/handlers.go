package server

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/shvz0/mmscrap/mmscrappers"
	"github.com/shvz0/mmscrap/stylometry"
	"gorm.io/gorm"
)

type MainPageHandler struct {
}

func (h MainPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var all []mmscrappers.Article

	alldb := Db.Scopes(Paginate(r)).Order("date desc").Find(&all)
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

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type DeltaHandler struct{}

func (h DeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stylometryRequest(w, r, stylometry.DeltaType)
}

type MendenhallHandler struct{}

func (h MendenhallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stylometryRequest(w, r, stylometry.MendenhallType)
}

type ChiSquredHandler struct{}

func (h ChiSquredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stylometryRequest(w, r, stylometry.ChiSquaredType)
}

func stylometryRequest(w http.ResponseWriter, r *http.Request, t stylometry.StylometryType) {

	var all []mmscrappers.Article

	Db.Find(&all)

	var corpuses []*stylometry.Corpus

	for _, v := range all {
		corpus := stylometry.NewCorpus(v.Text, v.Author)
		corpuses = append(corpuses, &corpus)
	}

	txt := r.FormValue("text")

	var res []stylometry.StylometryResult

	switch t {
	case stylometry.DeltaType:
		res = stylometry.DeltaMethod(corpuses, txt)
	case stylometry.ChiSquaredType:
		res = stylometry.ChiSquaredMethod(corpuses, txt)
	case stylometry.MendenhallType:
		res = stylometry.MendenhallMethod(corpuses, txt)
	}

	type payload struct {
		Message string
		Data    []stylometry.StylometryResult
	}

	p := payload{
		Message: "ok",
		Data:    res,
	}

	json, err := json.Marshal(p)

	if err != nil {
		log.Println(err)
		responseServerError(w, err)
		return
	}

	w.Write(json)
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
