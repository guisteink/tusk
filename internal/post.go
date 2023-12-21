package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Revision  string             `json:"revision"`
	Tips      string             `json:"tips"`
	CreatedAt time.Time          `json:"createdAt"`
	Tags      []string           `json:"tags"`
}
