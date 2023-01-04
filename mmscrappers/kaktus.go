package mmscrappers

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Kaktus struct {
	Headers map[string][]string
	reqDate time.Time
	wg      *sync.WaitGroup
}

func NewKaktus() Kaktus {
	return Kaktus{Headers: defaultHeaders(), wg: &sync.WaitGroup{}}
}

func (k Kaktus) ArticleListToday() ([]Article, error) {
	return k.ArticleListByTime(time.Now())
}

func (k Kaktus) ArticleListByTime(t time.Time) ([]Article, error) {

	var articles []Article

	k.reqDate = t

	articleListPage, err := GetPage(
		fmt.Sprintf("https://kaktus.media/?lable=8&date=%s&order=time",
			t.Format(DefaultDateFormat)),
		k.Headers)

	if err != nil {
		return nil, err
	}

	articleListItems, err := k.parseArticleListPage(articleListPage)

	if err != nil {
		return nil, err
	}

	ch := make(chan Article, len(articleListItems))

	for _, item := range articleListItems {
		k.wg.Add(1)
		go k.articlePageParse(item, ch)
	}

	k.wg.Wait()

	close(ch)

	for art := range ch {
		articles = append(articles, art)
	}

	return articles, nil
}

func (k Kaktus) parseArticleListPage(pageTxt io.Reader) (result []ArticleListItem, err error) {

	doc, err := htmlquery.Parse(pageTxt)

	if err != nil {
		return nil, err
	}

	namesNodes, err := htmlquery.QueryAll(doc, "//a[contains(@class, 'ArticleItem--name')]")
	if err != nil {
		return nil, err
	}

	imgNodes, err := htmlquery.QueryAll(doc, "//img[contains(@class, 'ArticleItem--image-img')]")
	if err != nil {
		return nil, err
	}

	timeNodes, err := htmlquery.QueryAll(doc, "//div[contains(@class, 'ArticleItem--time')]")
	if err != nil {
		return nil, err
	}

	for i, nn := range namesNodes {

		var imgSrc string
		var href string

		for _, attr := range nn.Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}

		for _, attr := range imgNodes[i].Attr {
			if attr.Key == "src" {
				imgSrc = attr.Val
			}
		}

		timeBuf := &bytes.Buffer{}

		collectText(timeNodes[i], timeBuf)

		timeStr := timeBuf.String()
		timeStr = strings.Trim(timeStr, "\n \t\r")

		timeHour, err := strconv.Atoi(timeStr[0:2])

		if err != nil {
			return nil, err
		}

		timeMin, err := strconv.Atoi(timeStr[3:])
		if err != nil {
			return nil, err
		}

		titleBuf := &bytes.Buffer{}

		collectText(nn, titleBuf)

		result = append(result,
			ArticleListItem{
				Url:      href,
				ImageUrl: imgSrc,
				Title:    titleBuf.String(),
				Date: time.Date(
					k.reqDate.Year(),
					k.reqDate.Month(),
					k.reqDate.Day(),
					timeHour,
					timeMin,
					0, 0, time.Local),
			})
	}

	return result, nil
}

func (k Kaktus) getAuthor(docN *html.Node) (string, error) {
	authorN, err := htmlquery.Query(docN, "//a[contains(@class, 'Article--author')]")
	if err != nil {
		return "", err
	}

	author := getNodeTxt(authorN)

	author = strings.ToLower(author)
	author = strings.Trim(author, " \n\t\r\x00")

	return author, nil
}

func (k Kaktus) articlePageParse(ali ArticleListItem, ch chan Article) (*Article, error) {

	articlePage, err := GetPage(ali.Url, k.Headers)

	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(articlePage)

	if err != nil {
		return nil, err
	}

	txtNodes, err := htmlquery.QueryAll(doc, "//div[contains(@class, 'Article--text')]//p")

	if err != nil {
		return nil, err
	}

	var txtBuf bytes.Buffer

	for _, tn := range txtNodes {
		collectText(tn, &txtBuf)
	}

	author, err := k.getAuthor(doc)

	if err != nil {
		author = ""
	}

	art := Article{
		Title:         ali.Title,
		Url:           ali.Url,
		Author:        author,
		Text:          txtBuf.String(),
		Date:          ali.Date,
		MassMediaName: "kaktus_media",
		ImgUrl:        ali.ImageUrl,
	}

	ch <- art

	k.wg.Done()

	return &art, nil

}
