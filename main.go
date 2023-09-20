package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
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

var wg sync.WaitGroup
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
			startTime := time.Now()
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
			endTime := time.Now()
			elapsedTime := endTime.Sub(startTime)
			utils.Log(fmt.Sprintf("Atualização de Posts Concluída! Levou %s", elapsedTime.String()))
			utils.Log(fmt.Sprintf("Próxima Atualização em %s", interval.String()))
			time.Sleep(interval)
		}
	default:
		for {
			startTime := time.Now()
			extractors.ExtractPagesOfForum(forum, &wg, client)
			wg.Wait()
			endTime := time.Now()
			elapsedTime := endTime.Sub(startTime)
			utils.Log(fmt.Sprintf("Extração Concluída! Levou %s", elapsedTime.String()))
			utils.Log(fmt.Sprintf("Próxima Extração em %s", interval.String()))
			time.Sleep(interval)
		}
	}
}
