package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/operations"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

var wg sync.WaitGroup

func parsePost(post entity.Post) {
	defer wg.Done()
	url := entity.URLForumPost("brdev", post)
	browser.Browser(url, func(p *colly.HTMLElement) {
		fullText := crawlers.PostCrawler(p)
		paragraphs := strings.Split(fullText, "\n")
		normalizedTitle := parsers.NormalizeTitle(post.Title)
		operations.CreateTextFile(normalizedTitle, paragraphs)
		operations.CreateAudioFile(normalizedTitle, paragraphs)
	})
}

func main() {
	start := time.Now()
	browser.Browser(entity.URLForum("brdev"), func(h *colly.HTMLElement) {
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
