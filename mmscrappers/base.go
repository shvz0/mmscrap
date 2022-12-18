package mmscrappers

import (
	"time"
)

type Article struct {
	Title         string
	Url           string
	Author        string
	Text          string
	Date          time.Time
	MassMediaName string
	ImgUrl        string
}

type ArticleListItem struct {
	Title    string
	Url      string
	ImageUrl string
	Date     time.Time
}

type MMscrapper interface {
	ArticleListToday() ([]Article, error)
}

type GetByPage interface {
	ArticleListByPage(pageNum int16) []Article
}

type GetByTime interface {
	ArticleListByDate(time.Time) []Article
}

var DefaultDateFormat string = "2006-01-02"
var DefaultDatetimeFormat string = "2006-01-02 15:04:05"

func GetAllScrappers() []MMscrapper {
	return []MMscrapper{
		NewKaktus(),
		NewNews24(),
	}
}
