package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (a *App) routes() chi.Router {
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
		http.ServeFile(w, r, filepath.Join(a.Opts.WebPath, "/favicon.ico"))
	})

	// file server for static content from /assets
	addFileServer(router, "/assets", http.Dir(a.Opts.WebPath))

	return router
}

func (a *App) handleNewPad() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := a.Storage.GetNextCounter()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		hash := a.Unique.Encode(int64(cnt))
		http.Redirect(w, r, "/"+hash, http.StatusFound)
	}
}

func (a *App) handleGetPad() http.HandlerFunc {
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
			tplPath := filepath.Join(a.Opts.WebPath, "templates/main.html")
			tpl, err = template.New("main").ParseFiles(tplPath)
			log.Printf("[DEBUG] parsed template: %s", tplPath)
		})
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		padname := chi.URLParam(r, "padname")
		content, err := a.Storage.GetPad(padname)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		tpl.Execute(w, data{
			Padname: padname,
			Content: content,
		})
	}
}

func (a *App) handleSetPad() http.HandlerFunc {
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

		err = a.Storage.SetPad(padname, content)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, &response{"error", padname})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, &response{"ok", padname})
	}
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
