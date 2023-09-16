package main

import (
	"strings"

	"github.com/gocolly/colly"
	browser "github.com/pedroosz/go-reddit-scrapper/src"
	"github.com/pedroosz/go-reddit-scrapper/src/crawlers"
	"github.com/pedroosz/go-reddit-scrapper/src/entity"
	"github.com/pedroosz/go-reddit-scrapper/src/operations"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
)

const BASE_URL = "https://old.reddit.com/"

func main() {
	browser.Browser(entity.URLForum("brdev"), func(h *colly.HTMLElement) {
		posts := crawlers.PostsCrawler(h)
		for i := 0; i < len(posts); i++ {
			post := posts[i]
			url := entity.URLForumPost("brdev", post)
			browser.Browser(url, func(p *colly.HTMLElement) {
				fullText := crawlers.PostCrawler(p)
				paragraphs := strings.Split(fullText, "\n")
				normalizedTitle := parsers.NormalizeTitle(post.Title)
				operations.CreateTextFile(normalizedTitle, paragraphs)
				operations.CreateAudioFile(normalizedTitle, paragraphs)
			})
		}
	})
}
