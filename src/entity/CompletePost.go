package entity

import "time"

type CompletePost struct {
	Url          string
	RawText      string
	Title        string
	Text         string
	Up           int
	CreationDate time.Time
	Comments     []Comment
	Kind         string
}
