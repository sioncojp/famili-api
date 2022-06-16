package v1todos

import (
	"net/http"
)

// TodoRepository...interfaceを使うことでDIPを解決する。mockも作成できるようになる
type Handler interface {
	Ctx(next http.Handler) http.Handler
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
