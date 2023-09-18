package entity

import "time"

type CompletePost struct {
	RawText      string
	Title        string
	Text         string
	Up           int
	CreationDate time.Time
}
