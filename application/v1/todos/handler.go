package v1todos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sioncojp/famili-api/domain"
	"github.com/sioncojp/famili-api/domain/model"
	"github.com/sioncojp/famili-api/domain/repository"
	httpresponse "github.com/sioncojp/famili-api/utils/http_response"
)

const (
	ErrorMessageNotFound        = "todo_not_found"
	ErrorMessageInvalidProvided = "invalid_todo_provided"
	ErrorMessageMissingArgument = "missing_argument"
	ErrorValidation             = "missing_validation"
)

var cv = &domain.CustomValidator{}

// Service...
type handler struct {
	repo repository.TodoRepository
}

// NewService create a instance of this service
func NewHandler(repo repository.TodoRepository) Handler {
	return &handler{repo}
}

// Ctx...アクセスした際に、既存の情報を保管する
func (s *handler) Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var todo model.Todo
		var err error

		// IDが空じゃないならidをベースにクエリを叩く
		if todoId := chi.URLParam(r, "id"); todoId != "" {
			todo, err = s.repo.GetById(domain.Id(todoId))
			if err != nil {
				httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageNotFound, "")
				return
			}
		} else {
			httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageInvalidProvided, "")
			return
		}

		ctx := context.WithValue(r.Context(), "todo", &todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// List...todoを取得してhttpを返す
func (s *handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := s.repo.List()
	if err != nil {
		httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageInvalidProvided, "")
		return
	}

	httpresponse.OK(w, r, http.StatusOK, "todos", out)
}

// Create...todoを作成してhttpを返す
func (s *handler) Create(w http.ResponseWriter, r *http.Request) {
	result := &model.Todo{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		httpresponse.Error(w, r, http.StatusBadRequest, ErrorMessageMissingArgument, "")
		return
	}
	if err := cv.Validate(result); err != nil {
		httpresponse.Error(w, r, http.StatusBadRequest, ErrorValidation, fmt.Sprintf("%s", err))
		return
	}
	result.Completed = false

	if err := s.repo.Create(result); err != nil {
		httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageInvalidProvided, "")
		return
	}
	httpresponse.OK(w, r, http.StatusCreated, "", nil)
}

// Update...todoを更新してhttpを返す
func (s *handler) Update(w http.ResponseWriter, r *http.Request) {
	result := &model.Todo{}
	todo := r.Context().Value("todo").(*model.Todo)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		httpresponse.Error(w, r, http.StatusBadRequest, ErrorMessageMissingArgument, "")
		return
	}

	if err := cv.Validate(result); err != nil {
		httpresponse.Error(w, r, http.StatusBadRequest, ErrorValidation, fmt.Sprintf("%s", err))
		return
	}

	todo.Title = result.Title
	todo.Description = result.Description
	todo.Completed = result.Completed

	if err := s.repo.Update(todo); err != nil {
		httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageInvalidProvided, "")
		return
	}

	httpresponse.OK(w, r, http.StatusOK, "", nil)
}

// Delete...todoを削除してhttpを返す
func (s *handler) Delete(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(*model.Todo)
	defer r.Body.Close()

	if err := s.repo.Delete(todo); err != nil {
		httpresponse.Error(w, r, http.StatusNotFound, ErrorMessageInvalidProvided, "")
		return
	}

	httpresponse.OK(w, r, http.StatusOK, "", nil)
}
