package browser

import (
	"fmt"

	"github.com/gocolly/colly"
	reddit "github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

func Browser(forum string) {
	c := colly.NewCollector(
		colly.AllowedDomains(reddit.BASE_URL),
	)

	c.OnRequest(func(r *colly.Request) {
		utils.Log(fmt.Sprintf("Visiting URL: %s", r.URL))
	})

	c.OnResponse(func(r *colly.Response) {
		utils.Log(fmt.Sprintf("Response from %s HTTP STATUS CODE % d", r.Request.URL, r.StatusCode))
	})

	c.Visit(reddit.URLForum(forum))
}
