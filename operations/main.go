package operations

import (
	"fmt"
	"log"
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/pedroosz/go-reddit-scrapper/parsers"
	"github.com/pedroosz/go-reddit-scrapper/utils"
)

const BASE_FOLDER = "files/"
const AUDIOS_FOLDER = BASE_FOLDER + "audios/"
const TEXTS_FOLDER = BASE_FOLDER + "texts/"

func CreateTextFile(title string, paragraphs []string) {

	if _, err := os.Stat(TEXTS_FOLDER); os.IsNotExist(err) {
		os.Mkdir(TEXTS_FOLDER, os.ModePerm)
	}

	fileName := title + ".txt"

	f, err := os.Create(TEXTS_FOLDER + fileName)

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

	utils.Log(title, "Arquivo de texto criado com sucesso")
}

func CreateAudioFile(title string, paragraphs []string) {
	speech := htgotts.Speech{Folder: AUDIOS_FOLDER + title, Language: voices.Portuguese}
	smallerPhrases := parsers.SplitPhrases(paragraphs, 30)

	done := make(chan bool)

	for i := 0; i < len(smallerPhrases); i++ {
		go func(phrase string, index int) {
			speech.CreateSpeechFile(phrase, fmt.Sprint(index))

			utils.Log(title, "Iniciando criação de audios.")
			utils.Log(title, "Audio "+fmt.Sprint(index)+" criado")

			done <- true
		}(smallerPhrases[i], i)
	}

	for i := 0; i < len(smallerPhrases); i++ {
		<-done
	}
}
