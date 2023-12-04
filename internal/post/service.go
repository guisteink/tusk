package post

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/guisteink/tusk/internal"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")

var MaxTimeout = 10

type Service struct {
	Repository Repository
}

type CreateResponse struct {
	Post internal.Post `json:"post"`
}

func (p Service) Create(post internal.Post) (CreateResponse, int, error) {
	if post.Body == "" {
		return CreateResponse{}, http.StatusBadRequest, ErrPostBodyEmpty
	}

	if utf8.RuneCountInString(post.Body) > 140 {
		return CreateResponse{}, http.StatusBadRequest, ErrPostBodyExceedsLimit
	}

	result, err := p.Repository.Insert(post, MaxTimeout)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return CreateResponse{}, http.StatusNotFound, err
		}
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	createdPost := internal.Post{
		Username:  result.Username,
		Title:     result.Title,
		Body:      result.Body,
		CreatedAt: result.CreatedAt,
		Tags:      result.Tags,
	}

	response := CreateResponse{Post: createdPost}
	return response, http.StatusCreated, nil
}
