package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/operations"
	"github.com/pedroosz/go-reddit-scrapper/parsers"
)

const BASE_URL = "https://old.reddit.com/"

func scrapTargetPost(title string, url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
	)

	c.OnHTML(".expando", func(e *colly.HTMLElement) {
		allParagraphs := make([]string, 0)

		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			allParagraphs = append(allParagraphs, h.Text)
		})

		normalizedTitle := parsers.NormalizeTitle(title)

		operations.CreateTextFile(normalizedTitle, allParagraphs)
		operations.CreateAudioFile(normalizedTitle, allParagraphs)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Received", r.StatusCode, r.Request.URL)
	})

	c.Visit(url)
}

func main() {

}
