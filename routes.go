package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Pad is a pad data for template
type Pad struct {
	Name    string
	Content string
}

func (a *App) routes() chi.Router {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	static := filepath.Join(workDir, "static")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.NoCache)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Get("/", a.handleNewPad())
	router.Get("/{padname}", a.handleGetPad())
	router.Post("/{padname}", a.handleSetPad())

	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "User-agent: *\n")
	})

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(static, "/favicon.ico"))
	})

	// file server for static content from /static
	addFileServer(router, "/static", http.Dir(static))
	return router
}

func (a *App) handleNewPad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := a.BoltBackend.GetNextCounter()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		hash := a.HashID.Encode(int64(cnt))
		http.Redirect(w, r, "/"+hash, 301)
	}
}

func (a *App) handleGetPad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		padname := chi.URLParam(r, "padname")
		content, err := a.BoltBackend.GetPad(padname)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tpl := template.New("main")
		tpl, err = tpl.ParseFiles("templates/main.html")
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{
				"message": "Template templates/main.html was not found",
			})
			return
		}

		tpl.Execute(w, Pad{
			Name:    padname,
			Content: content,
		})
	}
}

func (a *App) handleSetPad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// serves static files from /web
func addFileServer(r chi.Router, path string, root http.FileSystem) {
	origPath := path
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// don't show dirs, just serve files
		if strings.HasSuffix(r.URL.Path, "/") && len(r.URL.Path) > 1 && r.URL.Path != (origPath+"/") {
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}))
}
