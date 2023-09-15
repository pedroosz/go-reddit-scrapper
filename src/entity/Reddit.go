package reddit

import "fmt"

const (
	BASE_URL = "old.reddit.com"
)

func URLForum(name string) string {
	return fmt.Sprintf("https://%s/r/%s/", BASE_URL, name)
}
