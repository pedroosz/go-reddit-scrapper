package crawlers

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

/*
Recebe uma página HTML e devolve um slice contendo os Posts de um determinado Forum no Reddit.
*/
func PostsCrawler(h *colly.HTMLElement) []entity.Post {
	utils.Log("Starting to extract posts")
	var posts []entity.Post
	h.ForEach("#siteTable", func(i int, siteTable *colly.HTMLElement) {
		siteTable.ForEach(".thing", func(i int, e *colly.HTMLElement) {
			post := entity.BuildPostFromElement(e)
			utils.Log(fmt.Sprintf("Post Extraido: %s - %s ", post.Title, post.Url))
			posts = append(posts, post)
		})
	})
	return posts
}
