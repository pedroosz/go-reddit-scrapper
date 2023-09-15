package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/operations"
	"github.com/pedroosz/go-reddit-scrapper/parsers"
	"github.com/pedroosz/go-reddit-scrapper/utils"
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

		utils.Log("Parágrafos adquiridos", fmt.Sprint(len(allParagraphs)))

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
	c := colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
	)

	c.OnHTML("#siteTable", func(e *colly.HTMLElement) {
		e.ForEach(".thing", func(i int, h *colly.HTMLElement) {
			title := h.ChildText("a.title")
			url := h.ChildAttr("a.title", "href")

			utils.Log("Post encontrado", title)

			scrapTargetPost(title, BASE_URL+url)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Received", r.StatusCode, r.Request.URL)
	})

	c.Visit("https://old.reddit.com/r/EuSouOBabaca/")
}
