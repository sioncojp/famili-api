package httpresponse

import (
	"net/http"

	"github.com/go-chi/render"
)

// HttpRespondOK...2xx ~ 3xxのときに返す
func OK(w http.ResponseWriter, r *http.Request, statusCode int, field string, value interface{}) {
	response := make(map[string]interface{})
	response["ok"] = true
	if field != "" {
		response[field] = value
	}
	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

// HttpRespondError...4xx < 5xxのときに返すエラー
func Error(w http.ResponseWriter, r *http.Request, statusCode int, errorMessage, warn string) {
	response := make(map[string]interface{})
	response["ok"] = false
	response["error"] = errorMessage
	if warn != "" {
		response["warn"] = warn
	}
	render.Status(r, statusCode)
	render.JSON(w, r, response)
}
