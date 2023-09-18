package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/operations"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

var wg sync.WaitGroup

func createFiles(paragraphs []string, title string) {
	operations.CreateTextFile(title, paragraphs)
	operations.CreateAudioFile(title, paragraphs)
}

func parsePost(post entity.Post) {
	defer wg.Done()
	url := entity.URLForumPost("EuSouOBabaca", post)
	browser.Browser(url, func(p *colly.HTMLElement) {
		completePost := crawlers.PostCrawler(p, post)
		paragraphs := strings.Split(completePost.RawText, "\n")
		createFiles(paragraphs, completePost.Title)
	})
}

func config() {
	err := godotenv.Load(".env")
	if err != nil {
		utils.Log("Não foi possível carregar o arquivo de configuração ambiente")
		os.Exit(-1)
	}
}

func main() {
	config()
	start := time.Now()
	browser.Browser(entity.URLForum("EuSouOBabaca"), func(h *colly.HTMLElement) {
		posts := crawlers.PostsCrawler(h)
		wg.Add(len(posts) - 1)
		for i := 0; i < len(posts); i++ {
			go parsePost(posts[i])
		}
	})
	wg.Wait()
	end := time.Now()
	elapsed := end.Sub(start)
	utils.Log(fmt.Sprintf("Script time: %f", elapsed.Seconds()))
}
