package crawlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

func PostCrawler(e *colly.HTMLElement, post entity.Post) entity.CompletePost {
	fullPost := e.DOM.Find(".expando")
	postDate := e.DOM.Find(".side .linkinfo .date time")
	likes := e.DOM.Find(".midcol .score.likes")
	date, err := time.Parse("02 Jan 2006", postDate.Text())
	if err != nil {
		utils.Log(fmt.Sprintf("Não foi possível parsear a data %s", postDate.Text()))
	}
	postText := fullPost.Text()
	qttLikes, err := strconv.Atoi(likes.Text())
	if err != nil {
		utils.Log(fmt.Sprintf("Não foi possível parsear a quantidade de likes: %s", likes.Text()))
		qttLikes = -1
	}
	return entity.CompletePost{
		RawText:      postText,
		CreationDate: date,
		Title:        parsers.NormalizeTitle(post.Title),
		Up:           qttLikes,
	}
}
