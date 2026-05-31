package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/dotzero/pad/app/handlers"
)

func (a *App) makeHTTPServer(address string, port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", address, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (a *App) routes() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.NoCache)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)

	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "User-agent: *\n")
	})

	favicon := faviconHandler(a.StaticPath)
	router.Get("/favicon.ico", favicon)
	router.Get("/favicon.png", favicon)

	// file server for static content from /assets
	fileServer(router, "/assets", http.Dir(a.StaticPath))

	router.Get("/", handlers.Redirect(a.Storage, a.HashEncoder))
	router.Get("/raw/{padname}", handlers.Raw(a.Storage))
	router.Get("/md/{padname}", handlers.Markdown(a.Storage, a.Templates.Lookup("main.html")))
	router.Get("/{padname}", handlers.Get(a.Storage, a.Templates.Lookup("main.html")))
	router.Post("/{padname}", handlers.Set(a.Storage))

	return router
}

func faviconHandler(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticPath, "favicon.png"))
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
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
