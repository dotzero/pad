package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, map[string]string{"error": err.Error()})
}
