package utils

import "log"

func Log(title string, message string) {
	log.Printf("%s\n%s", title, message)
}
