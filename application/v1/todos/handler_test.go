package v1todos

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sioncojp/famili-api/domain/model"
	"github.com/sioncojp/famili-api/utils"
)

type TestCase struct {
	name           string
	parameter      string
	httpStatusCode int
}

var (
	contextKey = "todo"
	url        = "/v1/todos"
	urlId      = "/v1/todos/1"
)

func TestTodoList(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		{
			"ok",
			"",
			http.StatusOK,
		},
	}

	data := []model.Todo{
		{
			Title:       "1",
			Description: "hoge",
			Completed:   true,
		},
		{
			Title:       "2",
			Description: "fuga",
			Completed:   false,
		},
	}

	m := new(MockTodoService)
	m.On("List").Return(data, nil)
	s := NewHandler(m)

	for _, v := range cases {
		t.Run(
			v.name,
			func(tt *testing.T) {
				tt.Parallel()
				r := httptest.NewRequest(http.MethodGet, url, nil)
				w := httptest.NewRecorder()
				s.List(w, r)

				resp := w.Result()
				assert.Equal(tt, v.httpStatusCode, resp.StatusCode)
			},
		)
	}
}

func TestTodoCreate(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		{
			"ok",
			`{"title":"1","description":"hoge"}`,
			http.StatusCreated,
		},
		{
			"title above max size",
			fmt.Sprintf(`{"title":"%s","description":"hoge"}`, utils.MakeRandomString(51)),
			http.StatusBadRequest,
		},
		{
			"title below min size",
			`{"title":"","description":"hoge"}`,
			http.StatusBadRequest,
		},
		{
			"description above max size",
			fmt.Sprintf(`{"title":"1","description":"%s"}`, utils.MakeRandomString(101)),
			http.StatusBadRequest,
		},
		{
			"description below min size",
			`{"title":"1","description":""}`,
			http.StatusBadRequest,
		},
	}

	data := &model.Todo{
		Title:       "1",
		Description: "hoge",
		Completed:   false,
	}

	m := new(MockTodoService)
	m.On("Create", data).Return(nil).Once()
	s := NewHandler(m)

	for _, v := range cases {
		t.Run(
			v.name,
			func(tt *testing.T) {
				json := strings.NewReader(v.parameter)
				r := httptest.NewRequest(http.MethodPost, url, json)
				w := httptest.NewRecorder()
				s.Create(w, r)

				resp := w.Result()
				assert.Equal(tt, v.httpStatusCode, resp.StatusCode)
			},
		)
	}
}

func TestTodoUpdate(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		{
			"ok",
			`{"title":"2","description":"fuga", "completed":true}`,
			http.StatusOK,
		},
		{
			"title above max size",
			fmt.Sprintf(`{"title":"%s","description":"hoge"}`, utils.MakeRandomString(51)),
			http.StatusBadRequest,
		},
		{
			"title below min size",
			`{"title":"","description":"hoge"}`,
			http.StatusBadRequest,
		},
		{
			"description above max size",
			fmt.Sprintf(`{"title":"1","description":"%s"}`, utils.MakeRandomString(101)),
			http.StatusBadRequest,
		},
		{
			"description below min size",
			`{"title":"1","description":""}`,
			http.StatusBadRequest,
		},
	}

	data := &model.Todo{
		Title:       "1",
		Description: "hoge",
		Completed:   false,
	}

	dataUpdate := &model.Todo{
		Title:       "2",
		Description: "fuga",
		Completed:   true,
	}

	m := new(MockTodoService)
	m.On("Update", dataUpdate).Return(nil).Once()
	s := NewHandler(m)

	for _, v := range cases {
		t.Run(
			v.name,
			func(tt *testing.T) {
				json := strings.NewReader(v.parameter)
				r := httptest.NewRequest(http.MethodPut, urlId, json)
				ctx := context.WithValue(r.Context(), contextKey, data)
				w := httptest.NewRecorder()
				s.Update(w, r.WithContext(ctx))

				resp := w.Result()
				assert.Equal(tt, v.httpStatusCode, resp.StatusCode)
			},
		)
	}
}

func TestTodoDelete(t *testing.T) {
	t.Parallel()
	cases := []TestCase{
		{
			"ok",
			"",
			http.StatusOK,
		},
	}

	data := &model.Todo{
		Title:       "1",
		Description: "hoge",
		Completed:   false,
	}

	m := new(MockTodoService)
	m.On("Delete", data).Return(nil).Once()
	s := NewHandler(m)

	for _, v := range cases {
		t.Run(
			v.name,
			func(tt *testing.T) {
				r := httptest.NewRequest(http.MethodDelete, urlId, nil)
				ctx := context.WithValue(r.Context(), contextKey, data)
				w := httptest.NewRecorder()
				s.Delete(w, r.WithContext(ctx))

				resp := w.Result()
				assert.Equal(tt, v.httpStatusCode, resp.StatusCode)
			},
		)
	}
}
