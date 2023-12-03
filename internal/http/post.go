package internal

import (
	"time"

	"github.com/google/uuid"
)

func (p *Post) FormattedCreatedAt() string {
	return p.CreatedAt.Format("2006-01-02 15:04:05")
}

type Post struct {
	ID        uuid.UUID `json:"-"`
	Username  string    `json:"username"`
	Title 		string 		`json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}
