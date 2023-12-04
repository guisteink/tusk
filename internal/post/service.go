package post

import (
	"errors"
	"unicode/utf8"
	"net/http"

	"github.com/google/uuid"

	"github.com/guisteink/tusk/internal"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")

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

	postID, err := p.Repository.Insert(post)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return CreateResponse{}, http.StatusNotFound, err
		}
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	id, err := uuid.Parse(postID)
	if err != nil {
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	createdPost := internal.Post{
		ID:        id,
		Username:  post.Username,
		Title:     post.Title,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		Tags:      post.Tags,
	}

	response := CreateResponse{Post: createdPost}
	return response, http.StatusCreated, nil
}
