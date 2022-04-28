package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var errNotFound = errors.New("Pad is not found")

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case errNotFound:
		render.Status(r, http.StatusNotFound)
	default:
		render.Status(r, http.StatusInternalServerError)
	}

	render.JSON(w, r, map[string]string{"error": err.Error()})
}
