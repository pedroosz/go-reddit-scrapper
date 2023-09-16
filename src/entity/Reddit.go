package entity

import "fmt"

const (
	BASE_URL = "old.reddit.com"
)

/*
Constroi uma URL no formato esperado para acessar um forum do Reddit.
*/
func URLForum(name string) string {
	return fmt.Sprintf("https://%s/r/%s", BASE_URL, name)
}
