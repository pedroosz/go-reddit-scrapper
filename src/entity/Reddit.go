package entity

import (
	"fmt"
	"regexp"
)

const (
	BASE_URL = "old.reddit.com"
)

type Comment struct {
	Text    string
	Replies []Comment
}

/*
Constroi uma URL no formato esperado para acessar um forum do Reddit.
*/
func URLForum(name string) string {
	return fmt.Sprintf("https://%s/r/%s", BASE_URL, name)
}

func URLBaseReddit(route string) string {
	return fmt.Sprintf("https://%s%s", BASE_URL, route)
}

func ExtractCodeFromURL(url string) (string, error) {
	regexStr := fmt.Sprintf("https:\\/\\/%s\\/r\\/%s\\/comments\\/([a-z0-9]+)\\/", BASE_URL, "EuSouOBabaca")
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(url)
	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", fmt.Errorf("não foi possível extrair o id do post da URL %s", url)
}
