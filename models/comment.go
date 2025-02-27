package models

import "time"

type Comment struct {
	Id         int
	PostId     int
	AuthorId   int
	AuthorName string
	Text       string
	CreatedAt  time.Time
}
