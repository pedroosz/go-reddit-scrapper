package entity

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Post struct {
	Title string
	Url   string
}

/*
Recebe um elemento HTML e converte ele para um objeto Post.
*/
func BuildPostFromElement(e *colly.HTMLElement) Post {
	post := Post{
		Title: e.ChildText("a.title"),
		Url:   e.ChildAttr("a.title", "href"),
	}
	return post
}

/*
Recebe o nome do forum e o Post do Reddit para construir URL para acessar.
*/
func URLForumPost(forum string, post Post) string {
	return fmt.Sprintf("https://%s", BASE_URL) + post.Url
}
