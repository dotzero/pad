package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dotzero/pad/handlers"
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

	router.Get("/", handlers.Redirect(a.Storage, a.HashEncoder))
	router.Get("/{padname}", handlers.Get(a.Storage, filepath.Join(a.Opts.WebPath, "templates/main.html")))
	router.Post("/{padname}", handlers.Set(a.Storage))

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

// serves static files from /web
func addFileServer(r chi.Router, path string, root http.FileSystem) {
	origPath := path
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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
