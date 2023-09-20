package entity

type AWSVoice struct {
	Name           string
	WordsPerSecond float64
}

var Thiago = AWSVoice{
	Name:           "Thiago",
	WordsPerSecond: 3.1,
}
