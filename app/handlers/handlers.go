package handlers

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const (
	urlPad          = "padname"
	maxFormBodySize = 1 << 20 // 1 MiB
)

// Redirect handle redirects to new pads
func Redirect(s storage, e encoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := s.NextCounter()
		if err != nil {
			renderError(w, r, err)
			return
		}

		hash := e.Encode(int64(cnt)) //nolint:gosec

		http.Redirect(w, r, "/"+hash, http.StatusFound)
	}
}

// Get handle get specific pad
func Get(s storage, t tpl) http.HandlerFunc {
	return pad(s, t, "edit")
}

// Markdown handle get specific pad rendered as Markdown
func Markdown(s storage, t tpl) http.HandlerFunc {
	return pad(s, t, "markdown")
}

func pad(s storage, t tpl, mode string) http.HandlerFunc {
	type data struct {
		Padname  string
		Content  string
		Markdown template.HTML
		Mode     string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, urlPad)

		content, err := s.Get(padname)
		if err != nil {
			renderError(w, r, err)
			return
		}

		markdown, err := renderMarkdown(content)
		if err != nil {
			renderError(w, r, err)
			return
		}

		err = t.Execute(w, data{
			Padname:  padname,
			Content:  content,
			Markdown: markdown,
			Mode:     mode,
		})
		if err != nil {
			renderError(w, r, err)
		}
	}
}

func renderMarkdown(content string) (template.HTML, error) {
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	if err := markdown.Convert([]byte(content), &buf); err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil //nolint:gosec
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
		r.Body = http.MaxBytesReader(w, r.Body, maxFormBodySize)

		if err := r.ParseForm(); err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				render.Status(r, http.StatusRequestEntityTooLarge)
				render.JSON(w, r, map[string]string{"error": err.Error()})

				return
			}

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
