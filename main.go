package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/database"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/extractors"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var minutes, forum = getParams()
var interval = time.Duration(minutes) * time.Minute

func getForum() string {
	forum := os.Getenv("FORUM")
	if forum == "" {
		utils.Fatal("Forum deve ser forencedio", nil)
	}
	return forum
}

func getMinutes() int {
	minutes := os.Getenv("INTERVAL")
	if minutes == "" {
		utils.Fatal("Minutes deve ser fornecido", nil)
	}
	intMinutes, err := strconv.Atoi(minutes)

	if err != nil {
		utils.Fatal("Erro ao parsear minutos para inteiro", err)
	}
	return intMinutes
}

func getParams() (int, string) {
	forum := getForum()
	intMinutes := getMinutes()
	return intMinutes, forum
}

func main() {
	client = database.PrepareDatabase()
	switch os.Getenv("mode") {
	case "update":
		for {
			database.MapPostsOnDatabase(client, func(post *entity.CompletePost) {
				browser.Browser(post.Url, func(p *colly.HTMLElement) {
					completePost := crawlers.PostCrawler(p, entity.Post{
						Title: post.Title,
						Url:   post.Url,
					})
					completePost.Url = post.Url
					completePost.Comments = extractors.ExtratCommentsFromPost(post.Url)
					err := database.UpdatePost(post, &completePost, client)
					if err != nil {
						utils.Fatal("Erro ao atualizar registro do banco de dados", err)
					}
				})
			})
			time.Sleep(interval)
		}
	default:
		for {
			extractors.ExtractPagesOfForum(forum, client)
			time.Sleep(interval)
		}
	}
}
