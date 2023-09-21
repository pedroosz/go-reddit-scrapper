package crawlers

import (
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
)

func PostCrawler(e *colly.HTMLElement, post entity.Post) entity.CompletePost {
	fullPost := e.DOM.Find(".expando")
	postDate := e.DOM.Find(".side .linkinfo .date time")
	likes := e.DOM.Find(".midcol .score.likes")
	date, err := time.Parse("02 Jan 2006", postDate.Text())
	if err != nil {
		date = time.Now()
	}
	postText := fullPost.Text()
	qttLikes, err := strconv.Atoi(likes.Text())
	if err != nil {
		qttLikes = -1
	}
	return entity.CompletePost{
		RawText:      postText,
		CreationDate: date,
		Title:        post.Title,
		Up:           qttLikes,
	}
}
