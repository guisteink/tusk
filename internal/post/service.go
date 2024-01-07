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

var MaxCharactersPerPost = 20000

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

	if utf8.RuneCountInString(post.Body) > MaxCharactersPerPost {
		return CreateResponse{}, http.StatusBadRequest, ErrPostBodyExceedsLimit
	}

	result, err := p.Repository.Insert(post)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return CreateResponse{}, http.StatusNotFound, err
		}
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	createdPost := internal.Post{
		ID:        result.ID,
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
			return []internal.Post{}, http.StatusOK, nil
		}
		return nil, http.StatusInternalServerError, err
	}

	if len(posts) == 0 {
		return []internal.Post{}, http.StatusOK, nil
	}

	return posts, http.StatusOK, nil
}

func (s Service) DeleteByID(param string) (CreateResponse, int, error) {
	if param == "" {
		return CreateResponse{}, http.StatusBadRequest, ErrIdEmpty
	}

	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return CreateResponse{}, http.StatusBadRequest, fmt.Errorf("invalid id format: %v", err)
	}

	deletedPost, err := s.Repository.Delete(id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return CreateResponse{}, http.StatusNotFound, err
		}
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	response := CreateResponse{Post: deletedPost}
	return response, http.StatusOK, nil
}

func (s Service) UpdateByID(id string, updatedPost internal.Post) (CreateResponse, int, error) {
	if updatedPost.Body == "" {
		return CreateResponse{}, http.StatusBadRequest, ErrPostBodyEmpty
	}

	if utf8.RuneCountInString(updatedPost.Body) > MaxCharactersPerPost {
		return CreateResponse{}, http.StatusBadRequest, ErrPostBodyExceedsLimit
	}

	if id == "" {
		return CreateResponse{}, http.StatusBadRequest, ErrIdEmpty
	}

	// Convertendo ObjectID para string
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return CreateResponse{}, http.StatusBadRequest, fmt.Errorf("invalid id format: %v", err)
	}

	result, err := s.Repository.Update(oid, updatedPost)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return CreateResponse{}, http.StatusNotFound, err
		}
		return CreateResponse{}, http.StatusInternalServerError, err
	}

	responsePost := internal.Post{
		Username:  result.Username,
		Title:     result.Title,
		Body:      result.Body,
		CreatedAt: result.CreatedAt,
		Tags:      result.Tags,
	}

	response := CreateResponse{Post: responsePost}
	return response, http.StatusOK, nil
}

func (s Service) FindWithPagination(startIndex, endIndex int) ([]internal.Post, int, error) {
	posts, totalItems, err := s.Repository.FindWithPagination(startIndex, endIndex)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			return []internal.Post{}, 0, nil
		}
		return nil, 0, err
	}

	if len(posts) == 0 {
		return []internal.Post{}, 0, nil
	}

	return posts, totalItems, nil
}
