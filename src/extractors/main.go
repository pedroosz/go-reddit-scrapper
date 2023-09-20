package extractors

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/database"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExtractPost(post entity.Post, wg *sync.WaitGroup, client *mongo.Client) {
	defer wg.Done()
	url := entity.URLForumPost(os.Getenv("FORUM"), post)
	post.Url = url
	if database.PostExistsOnCollection(post, client) {
		utils.Log(fmt.Sprintf("Post (%s) já foi parseado", post.Url))
		return
	}
	utils.Log(fmt.Sprintf("Post (%s) não foi parseado", post.Url))
	browser.Browser(url, func(p *colly.HTMLElement) {
		completePost := crawlers.PostCrawler(p, post)
		completePost.Url = url
		completePost.Comments = ExtratCommentsFromPost(url)
		parsers.ParseKind(&completePost)
		database.InsertPostsOnCollection(completePost, client)
	})
}

func ExtractPagesOfForum(forum string, wg *sync.WaitGroup, client *mongo.Client) {
	browser.Browser(entity.URLForum(forum), func(h *colly.HTMLElement) {
		posts := crawlers.PostsCrawler(h)
		wg.Add(len(posts))
		for i := 0; i < len(posts); i++ {
			go ExtractPost(posts[i], wg, client)
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
	utils.Log("Extraindo comentários de " + urlAll)
	browser.Browser(urlAll, func(h *colly.HTMLElement) {
		commentArea := h.DOM.Find(".commentarea")
		comments = ExtractCommentsFromContainer(commentArea)
	})
	utils.Log("Comentários de " + urlAll + " extraídos")
	return comments
}

func ExtratCommentsFromPost(url string) []entity.Comment {
	var comments []entity.Comment
	browser.Browser(url, func(h *colly.HTMLElement) {
		linkToSeeAllComments, exists := h.DOM.Find(".panestack-title .title-button").Attr("href")
		if !exists {
			utils.Log("Link para comentários não existe - Realizando extração na página atual")
			comments = ExtractAllComments(url)
			return
		}
		utils.Log("Link para comentários existe! Começando a realizar extração na página com todos os comentários")
		comments = ExtractAllComments(entity.URLBaseReddit(linkToSeeAllComments))
	})
	return comments
}
