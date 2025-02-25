package models

import "time"

type Article struct {
	Id        int
	AuthorId  int
	Title     string
	Body      string
	IsPosted  bool
	CoverURL  string
	CreatedAt time.Time
}
