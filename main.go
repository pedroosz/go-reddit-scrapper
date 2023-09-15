package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pedroosz/go-reddit-scrapper/audios"
	"github.com/pedroosz/go-reddit-scrapper/texts"
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

		fmt.Printf("=====================\n")
		fmt.Printf("SCRAPPING CONCLUÍDO\n")
		fmt.Printf("TITULO: %s\nURL: %s\nPARAGRAFOS: %d", title, url, len(allParagraphs))
		fmt.Printf("=====================\n")

		titleWithoutSpaces := strings.ReplaceAll(title, " ", "_")
		titleWithoutSpacesAndSlashes := strings.ReplaceAll(titleWithoutSpaces, "/", "_")

		texts.CreateFiles(titleWithoutSpacesAndSlashes, allParagraphs)
		audios.CreateAudio(titleWithoutSpacesAndSlashes, allParagraphs)
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

			fmt.Printf("=====================\n")
			fmt.Printf("POST ENCONTRADO")
			fmt.Printf("TITULO: %s\nURL: %s", title, BASE_URL+url)
			fmt.Printf("=====================\n")

			scrapTargetPost(title, BASE_URL+url)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Received", r.StatusCode, r.Request.URL)
	})

	c.Visit("https://old.reddit.com/r/brdev/")
}
