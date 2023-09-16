package crawlers

import (
	"github.com/gocolly/colly"
)

func PostCrawler(e *colly.HTMLElement) string {
	fullPost := e.DOM.Find(".expando")
	postText := fullPost.Text()
	return postText
}
