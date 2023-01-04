package mmscrappers

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type News24 struct {
	Headers  map[string][]string
	reqDate  time.Time
	wg       *sync.WaitGroup
	throttle time.Duration
}

func NewNews24() News24 {
	return News24{
		Headers:  defaultHeaders(),
		wg:       &sync.WaitGroup{},
		throttle: 300 * time.Millisecond}
}

func (n News24) GetURL() string {
	return "https://24.kg/"
}

func (n News24) ArticleListByPage(pageNum int16) {

}

func (n News24) ArticleListToday() ([]Article, error) {

	var articles []Article

	articleListPage, err := GetPage(n.GetURL(), defaultHeaders())

	if err != nil {
		return nil, err
	}

	articleListItems, err := n.parseArticleListPageOneDay(articleListPage)

	if err != nil {
		return nil, err
	}

	ch := make(chan Article, len(articleListItems))

	for _, item := range articleListItems {
		n.wg.Add(1)
		time.Sleep(n.throttle)
		go n.articlePageParse(item, ch)
	}

	n.wg.Wait()

	close(ch)

	for art := range ch {
		articles = append(articles, art)
	}

	return articles, nil
}

func (n News24) parseArticleListPageOneDay(page io.Reader) ([]ArticleListItem, error) {

	var artList []ArticleListItem

	doc, err := htmlquery.Parse(page)

	if err != nil {
		return nil, err
	}

	listStartNode, err := htmlquery.Query(doc, "//div[contains(@class, 'lineDate')]")

	if err != nil {
		return nil, err
	}

	if listStartNode == nil {
		return make([]ArticleListItem, 0), errors.New("List is empty")
	}

	for curr := listStartNode.NextSibling; curr != nil; curr = curr.NextSibling {
		if strings.Contains(getAttr("class", *curr), "lineDate") {
			break
		}

		if strings.Contains(getAttr("class", *curr), "one") {
			article, err := n.articleParseListItem(*curr)
			if err != nil {
				return nil, err
			}
			artList = append(artList, article)
		}
	}

	return artList, nil
}

func (n News24) articleParseListItem(node html.Node) (ArticleListItem, error) {

	titleNode := node.FirstChild.NextSibling.NextSibling.NextSibling
	hrefNode := titleNode.FirstChild.NextSibling

	url := n.GetURL() + getAttr("href", *hrefNode)[1:]
	title := strings.Trim(getNodeTxt(titleNode), " \n\t\r")

	return ArticleListItem{Url: url, Title: title}, nil
}

func (n News24) articlePageParse(i ArticleListItem, ch chan Article) (Article, error) {

	var a Article

	page, err := GetPage(i.Url, defaultHeaders())

	if err != nil {
		return Article{}, err
	}

	doc, err := htmlquery.Parse(page)

	if err != nil {
		return Article{}, err
	}

	node, err := htmlquery.Query(doc, "//div[@itemprop=\"articleBody\"]")

	if err != nil {
		return Article{}, err
	}

	a.Text = getNodeTxt(node)
	a.Url = i.Url
	a.Title = i.Title
	a.MassMediaName = "news_24kg"
	a.Date, err = n.articlePageParseDate(doc)
	a.Author = n.getAuthor(doc)

	if err != nil {
		return Article{}, err
	}

	ch <- a

	n.wg.Done()

	return a, nil
}

func (n News24) getAuthor(docN *html.Node) string {

	tn, _ := htmlquery.Query(docN, "//span[@itemprop='author']")

	tStr := getNodeTxt(tn)

	tStr = strings.ToLower(tStr)
	tStr = strings.Trim(tStr, " \n\t\r\x00")

	return tStr
}

func (n News24) articlePageParseDate(docN *html.Node) (time.Time, error) {

	tn, err := htmlquery.Query(docN, "//span[@itemprop='datePublished']")

	tStr := getNodeTxt(tn)

	t, err := n.parseDateRu(tStr)

	if err != nil {
		return time.Time{}, err
	}

	return *t, nil
}

func (n News24) parseDateRu(dateStr string) (*time.Time, error) {

	dateLowerTrimmed := strings.Trim(strings.ToLower(dateStr), " \n\r\t")

	r, err := regexp.Compile(`(?is)(\d{2}):(\d{2}),\s+(\d{2})\s+([А-Я]+)\s+(\d{4})`)

	m := r.FindSubmatch([]byte(dateLowerTrimmed))

	hour, err := strconv.Atoi(string(m[1]))
	min, err := strconv.Atoi(string(m[2]))
	day, err := strconv.Atoi(string(m[3]))
	month, err := n.monthByStrRu(string(m[4]))
	year, err := strconv.Atoi(string(m[5]))

	if err != nil {
		return nil, err
	}

	t := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local)

	return &t, nil
}

func (n News24) monthByStrRu(s string) (time.Month, error) {

	if strings.Contains(s, "янв") {
		return time.January, nil
	}
	if strings.Contains(s, "фев") {
		return time.February, nil
	}
	if strings.Contains(s, "март") {
		return time.March, nil
	}
	if strings.Contains(s, "апр") {
		return time.April, nil
	}
	if strings.Contains(s, "май") {
		return time.May, nil
	}
	if strings.Contains(s, "июн") {
		return time.June, nil
	}
	if strings.Contains(s, "июл") {
		return time.July, nil
	}
	if strings.Contains(s, "авг") {
		return time.August, nil
	}
	if strings.Contains(s, "сент") {
		return time.September, nil
	}
	if strings.Contains(s, "окт") {
		return time.October, nil
	}
	if strings.Contains(s, "ноя") {
		return time.November, nil
	}
	if strings.Contains(s, "дек") {
		return time.December, nil
	}

	return 0, fmt.Errorf("String does not match any month")
}
