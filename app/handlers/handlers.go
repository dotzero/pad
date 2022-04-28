package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	urlPad = "padname"
)

// Redirect handle redirects to new pads
func Redirect(s storage, e encoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := s.NextCounter()
		if err != nil {
			renderError(w, r, err)
			return
		}

		hash := e.Encode(int64(cnt))

		http.Redirect(w, r, "/"+hash, http.StatusFound)
	}
}

// Get handle get specific pad
func Get(s storage, t tpl) http.HandlerFunc {
	type data struct {
		Padname string
		Content string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, urlPad)

		content, err := s.Get(padname)
		if err != nil {
			renderError(w, r, err)
			return
		}

		err = t.Execute(w, data{
			Padname: padname,
			Content: content,
		})
		if err != nil {
			renderError(w, r, err)
		}
	}
}

// Raw handle get specific pad in plain text
func Raw(s storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, urlPad)

		content, err := s.Get(padname)
		if err != nil {
			renderError(w, r, err)
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, content)
	}
}

// Set handle set specific pad
func Set(s storage) http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
		Padname string `json:"padname,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			renderError(w, r, err)
			return
		}

		padname := chi.URLParam(r, urlPad)
		content := r.Form.Get("t")

		if err := s.Set(padname, content); err != nil {
			renderError(w, r, err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &response{Message: "ok", Padname: padname})
	}
}
