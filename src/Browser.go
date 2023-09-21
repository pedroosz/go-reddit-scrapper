package browser

import (
	"github.com/gocolly/colly"
	reddit "github.com/pedroosz/go-reddit-scrapper/src/entity"
)

type callback func(h *colly.HTMLElement)

/*
Faz o papel de um navegador o web e realiza a requisição para obter o arquivo HTML do forum daquele reddit.
É passado dois argumentos: uma URL da página que você deseja realizar o scrapping e uma função callback. A função callback recebe como argumento
o elemento HTML inteiro da página.
*/
func Browser(url string, callback callback) {
	c := colly.NewCollector(
		colly.AllowedDomains(reddit.BASE_URL),
	)

	c.OnHTML("html", func(h *colly.HTMLElement) {
		callback(h)
	})

	c.Visit(url)
}
