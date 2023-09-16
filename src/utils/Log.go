package utils

import (
	"log"
)

/*
Realiza a configuração das flags do logger.
https://pkg.go.dev/log#pkg-constants
*/
func configureFlags() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}

/*
Configura e realiza o logging da mensagem.
*/
func Log(message string) {
	configureFlags()
	log.Println(message)
}
