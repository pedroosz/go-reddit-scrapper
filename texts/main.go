package texts

import (
	"fmt"
	"log"
	"os"
)

func CreateFiles(title string, paragraphs []string) {
	folder := "./files"

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}

	f, err := os.Create(folder + "/" + title + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for i := 0; i < len(paragraphs); i++ {
		_, err2 := f.WriteString(
			paragraphs[i] + "\n",
		)

		if err2 != nil {
			log.Fatal(err2)
		}
	}

	fmt.Println(title, "pronto.")
}
