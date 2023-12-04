package post

import (
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/guisteink/tusk/internal"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")
var ErrIdEmpty = errors.New("id empty")

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

func (s Service) FindByID(param string) (internal.Post, int, error) {
	if param == "" {
		return internal.Post{}, http.StatusBadRequest, ErrIdEmpty
	}

	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return internal.Post{}, http.StatusBadRequest, fmt.Errorf("invalid id format: %v", err)
	}

	posts, err := s.Repository.Find(primitive.M{"_id": id})
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return internal.Post{}, http.StatusNotFound, err
		}
		return internal.Post{}, http.StatusInternalServerError, err
	}

	if len(posts) == 0 {
		return internal.Post{}, http.StatusNotFound, ErrPostNotFound
	}

	foundPost := posts[0]
	return foundPost, http.StatusOK, nil
}

func (s Service) FindAll() ([]internal.Post, int, error) {
	posts, err := s.Repository.Find(bson.M{})
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return []internal.Post{}, http.StatusNotFound, err
		}
		return []internal.Post{}, http.StatusInternalServerError, err
	}

	if len(posts) == 0 {
		return []internal.Post{}, http.StatusNotFound, ErrPostNotFound
	}

	return posts, http.StatusOK, nil
}
