package parsers

import "strings"

func NormalizeTitle(title string) string {
	titleWithoutSpaces := strings.ReplaceAll(title, " ", "_")
	titleWithoutSlashes := strings.ReplaceAll(titleWithoutSpaces, "/", "_")
	titleWithoutQuestionMarks := strings.ReplaceAll(titleWithoutSlashes, "?", "")

	return titleWithoutQuestionMarks
}

func SplitPhrases(bigPhrases []string, maxWords int) []string {
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
