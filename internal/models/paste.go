package models

import "time"

type Paste struct {
	ID       int
	Content  string
	Language string
	Created  time.Time
	Expires  time.Time
}
