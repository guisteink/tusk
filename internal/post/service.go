package post

import (
	"fmt"
	"errors"
	"unicode/utf8"

	"github.com/guisteink/tusk/internal"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")

type Service struct {
	Repository Repository
}

func (p Service) Create(post internal.Post) error {
	fmt.Printf("\n\n[service create] = %+v\n", post, "\n\n")

	if post.Body == "" {
		return ErrPostBodyEmpty
	}

	if utf8.RuneCountInString(post.Body) > 140 {
		return ErrPostBodyExceedsLimit
	}

	return p.Repository.Insert(post)
}
