package utils

import (
	"log"
)

func configureFlags() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}

func Log(message string) {
	configureFlags()
	log.Println(message)
}
