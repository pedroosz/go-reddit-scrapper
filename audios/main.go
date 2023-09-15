package audios

import (
	"fmt"
	"log"
	"strings"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

func splitPhrases(bigPhrases []string, maxWords int) []string {
	var smallerPhrases []string

	for _, phrase := range bigPhrases {
		words := strings.Fields(phrase)
		for i := 0; i < len(words); i += maxWords {
			end := i + maxWords
			if end > len(words) {
				end = len(words)
			}
			smallerPhrases = append(smallerPhrases, strings.Join(words[i:end], " "))
		}
	}

	return smallerPhrases
}

func CreateAudio(title string, paragraphs []string) {
	speech := htgotts.Speech{Folder: "files/audios/" + title, Language: voices.Portuguese}
	smallerPhrases := splitPhrases(paragraphs, 30)

	// Create a channel to coordinate goroutines
	done := make(chan bool)

	for i := 0; i < len(smallerPhrases); i++ {
		go func(phrase string, index int) {
			speech.CreateSpeechFile(phrase, fmt.Sprint(index))
			log.Default().Println("Audio numero", index, "do post", title, "criado")
			done <- true
		}(smallerPhrases[i], i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < len(smallerPhrases); i++ {
		<-done
	}
}
