package crawlers

import (
	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
)

/*
Recebe uma página HTML e devolve um slice contendo os Posts de um determinado Forum no Reddit.
*/
func PostsCrawler(h *colly.HTMLElement) []entity.Post {
	var posts []entity.Post
	h.ForEach("#siteTable", func(i int, siteTable *colly.HTMLElement) {
		siteTable.ForEach(".thing", func(i int, e *colly.HTMLElement) {
			post := entity.BuildPostFromElement(e)
			posts = append(posts, post)
		})
	})
	return posts
}
