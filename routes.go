package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

// Pad is a pad data for template
type Pad struct {
	Name    string
	Content string
}

func (a *App) newPad(w http.ResponseWriter, r *http.Request) {
	cnt, err := a.BoltBackend.GetNextCounter()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hash := a.HashID.Encode(int64(cnt))
	http.Redirect(w, r, "/"+hash, 301)
}

func (a *App) getPad(w http.ResponseWriter, r *http.Request) {
	padname := chi.URLParam(r, "padname")
	content, err := a.BoltBackend.GetPad(padname)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tpl := template.New("main")
	tpl, _ = tpl.ParseFiles("templates/main.html")
	tpl.Execute(w, Pad{
		Name:    padname,
		Content: content,
	})
}

func (a *App) setPad(w http.ResponseWriter, r *http.Request) {
	var err error
	err = r.ParseForm()
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "error",
		})
		return
	}
	padname := chi.URLParam(r, "padname")
	content := r.Form.Get("t")

	err = a.BoltBackend.SetPad(padname, content)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"message": "error",
			"padname": padname,
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "ok",
		"padname": padname,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
