package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/zero-pkg/tpl"

	"github.com/dotzero/pad/handlers"
	"github.com/dotzero/pad/hash"
	"github.com/dotzero/pad/storage"
)

type hashEncoder interface {
	Encode(num int64) string
}

type padStorage interface {
	Get(name string) (value string, err error)
	Set(name string, value string) error
	NextCounter() (next uint64, err error)
}

// App is a Pad app
type App struct {
	Opts
	Storage     padStorage
	HashEncoder hashEncoder
	Templates   *tpl.Templates
}

// New prepares application and return it
func New(opts Opts) (*App, error) {
	if err := makeDirs(opts.DatabasePath); err != nil {
		return nil, err
	}

	store, err := storage.New(opts.DatabasePath, "pad.db")
	if err != nil {
		return nil, err
	}

	return &App{
		Opts:        opts,
		Storage:     store,
		HashEncoder: hash.New(opts.SecretKey, 3),
		Templates:   tpl.Must(tpl.New().ParseDir(filepath.Join(opts.AssetsPath, "templates"), ".html")),
	}, nil
}

// Run the listener
func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.Opts.Host, a.Opts.Port)
	log.Printf("[INFO] http server listen at: http://" + addr)

	router := a.routes()
	return http.ListenAndServe(addr, router)
}

func (a *App) routes() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.NoCache)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)

	router.Get("/", handlers.Redirect(a.Storage, a.HashEncoder))
	router.Get("/raw/{padname}", handlers.Raw(a.Storage))
	router.Get("/{padname}", handlers.Get(a.Storage, a.Templates))
	router.Post("/{padname}", handlers.Set(a.Storage))

	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "User-agent: *\n")
	})

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(a.Opts.AssetsPath, "/favicon.ico"))
	})

	// file server for static content from /assets
	addFileServer(router, "/assets", http.Dir(a.Opts.AssetsPath))

	return router
}

func makeDirs(dirs ...string) error {
	// exists returns whether the given file or directory exists or not
	exists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}

	for _, dir := range dirs {
		ex, err := exists(dir)
		if err != nil {
			return fmt.Errorf("can't check directory status for %s", dir)
		}
		if !ex {
			if e := os.MkdirAll(dir, 0700); e != nil {
				return fmt.Errorf("can't make directory %s", dir)
			}
		}
	}
	return nil
}

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
