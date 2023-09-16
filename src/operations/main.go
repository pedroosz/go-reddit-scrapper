package operations

import (
	"fmt"
	"log"
	"os"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/pedroosz/go-reddit-scrapper/src/parsers"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

const BASE_FOLDER = "files/"
const AUDIOS_FOLDER = BASE_FOLDER + "audios/"
const TEXTS_FOLDER = BASE_FOLDER + "texts/"

func createBaseFolder() {
	if _, err := os.Stat(BASE_FOLDER); os.IsNotExist(err) {
		os.Mkdir(BASE_FOLDER, os.ModePerm)
	}
}

func CreateTextFile(title string, paragraphs []string) {

	createBaseFolder()

	if _, err := os.Stat(TEXTS_FOLDER); os.IsNotExist(err) {
		os.Mkdir(TEXTS_FOLDER, os.ModePerm)
	}

	fileName := title + ".txt"

	f, err := os.Create(TEXTS_FOLDER + fileName)

	if err != nil {
		utils.Log(err.Error())
		os.Exit(-1)
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
}

func CreateAudioFile(title string, paragraphs []string) {

	createBaseFolder()

	speech := htgotts.Speech{Folder: AUDIOS_FOLDER + title, Language: voices.Portuguese}
	smallerPhrases := parsers.SplitPhrases(paragraphs, 30)

	done := make(chan bool)

	for i := 0; i < len(smallerPhrases); i++ {
		go func(phrase string, index int) {
			speech.CreateSpeechFile(phrase, fmt.Sprint(index))

			done <- true
		}(smallerPhrases[i], i)
	}

	for i := 0; i < len(smallerPhrases); i++ {
		<-done
	}
}
