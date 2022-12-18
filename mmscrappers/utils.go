package mmscrappers

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func PrintArticle(a Article) {
	fmt.Printf("Link:%s\nTitle:%s\nTXT:%s\n%s\n", a.Url, a.Title, a.Text, strings.Repeat("-", 50))
}

func SaveTodaysArticlesToDB(db *gorm.DB) {

	ss := GetAllScrappers()
	for _, s := range ss {
		articles, err := s.ArticleListToday()

		if err != nil {
			panic(err)
		}

		for _, a := range articles {
			if &a != nil {
				db.Create(&a)
			}

			PrintArticle(a)

			if err != nil {
				if strings.Contains(err.Error(), "duplicate") {
					fmt.Println("Already exists")
					continue
				} else {
					panic(err)
				}
			}
		}
	}
}
