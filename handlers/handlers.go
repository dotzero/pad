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
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
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
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		err = t.Execute(w, data{
			Padname: padname,
			Content: content,
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
		}
	}
}

// Raw handle get specific pad in plain text
func Raw(s storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, urlPad)
		content, err := s.Get(padname)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
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
		var err error
		err = r.ParseForm()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, &response{Message: "error"})
			return
		}

		padname := chi.URLParam(r, urlPad)
		content := r.Form.Get("t")

		err = s.Set(padname, content)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, &response{Message: "error", Padname: padname})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &response{Message: "ok", Padname: padname})
	}
}
