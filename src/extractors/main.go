package extractors

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/database"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExtractPost(post entity.Post, client *mongo.Client) {
	url := entity.URLForumPost(os.Getenv("FORUM"), post)
	post.Url = url
	if database.PostExistsOnCollection(post, client) {
		return
	}
	browser.Browser(url, func(p *colly.HTMLElement) {
		completePost := crawlers.PostCrawler(p, post)
		completePost.Url = url
		completePost.Comments = ExtratCommentsFromPost(url)
		parsers.ParseKind(&completePost)
		database.InsertPostsOnCollection(completePost, client)
	})
}

func ExtractPagesOfForum(forum string, client *mongo.Client) {
	browser.Browser(entity.URLForum(forum), func(h *colly.HTMLElement) {
		posts := crawlers.PostsCrawler(h)
		for i := 0; i < len(posts); i++ {
			ExtractPost(posts[i], client)
		}
	})
}

func ExtractCommentsFromContainer(containner *goquery.Selection) []entity.Comment {
	comments := make([]entity.Comment, 0)
	containner.Find(".sitetable").Children().Each(func(i int, s *goquery.Selection) {
		comment := entity.Comment{}
		commentText := s.Find(".usertext-body .md")
		comment.Text = commentText.First().Text()
		if len(strings.TrimSpace(comment.Text)) == 0 {
			return
		}
		child := s.Find(".child")
		if len(child.Find(".usertext-body .md").First().Text()) == 0 {
			comments = append(comments, comment)
			return
		}
		comment.Replies = ExtractCommentsFromContainer(child)
		comments = append(comments, comment)
	})
	return comments
}

func ExtractAllComments(urlAll string) []entity.Comment {
	var comments []entity.Comment
	browser.Browser(urlAll, func(h *colly.HTMLElement) {
		commentArea := h.DOM.Find(".commentarea")
		comments = ExtractCommentsFromContainer(commentArea)
	})
	return comments
}

func ExtratCommentsFromPost(url string) []entity.Comment {
	var comments []entity.Comment
	browser.Browser(url, func(h *colly.HTMLElement) {
		linkToSeeAllComments, exists := h.DOM.Find(".panestack-title .title-button").Attr("href")
		if !exists {
			comments = ExtractAllComments(url)
			return
		}
		comments = ExtractAllComments(entity.URLBaseReddit(linkToSeeAllComments))
	})
	return comments
}
