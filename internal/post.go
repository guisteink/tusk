package internal

import (
	"time"
)

type Post struct {
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      []string  `json:"tags"`
}
