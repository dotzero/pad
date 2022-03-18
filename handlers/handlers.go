package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type hashEncoder interface {
	Encode(num int64) string
}

type padStorage interface {
	Get(name string) (value string, err error)
	Set(name string, value string) error
	NextCounter() (next uint64, err error)
}

// Redirect handle redirects to new pads
func Redirect(s padStorage, h hashEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := s.NextCounter()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		hash := h.Encode(int64(cnt))
		http.Redirect(w, r, "/"+hash, http.StatusFound)
	}
}

// Get handle get specific pad
func Get(s padStorage, tplPath string) http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)

	type data struct {
		Padname string
		Content string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			name := filepath.Base(tplPath)
			tpl, err = template.New(name).ParseFiles(tplPath)
			log.Printf("[DEBUG] parsed template: %s", tplPath)
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		padname := chi.URLParam(r, "padname")
		content, err := s.Get(padname)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		err = tpl.Execute(w, data{
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
func Raw(s padStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, "padname")
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
func Set(s padStorage) http.HandlerFunc {
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

		padname := chi.URLParam(r, "padname")
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
