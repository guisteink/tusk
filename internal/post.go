package internal

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"-"`
	Username  string    `json:"username"`
	Title 		string 		`json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      []string  `json:"tags"`
}
