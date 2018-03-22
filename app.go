package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dotzero/pad/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// App is a Pad app
type App struct {
	Config      *Configuration
	Router      *chi.Mux
	BoltBackend *service.BoltBackend
	HashID      *service.HashID
}

// NewPadApp returns a client of the Pad app
func NewPadApp(cfg *Configuration) (*App, error) {
	boltBackend, err := service.NewBoltBackend(cfg.DB)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:      cfg,
		Router:      chi.NewRouter(),
		BoltBackend: boltBackend,
		HashID:      service.NewHashID(cfg.Salt, 3),
	}, nil
}

// Initialize adds middlewares, routes and static handlers to app
func (a *App) Initialize(cfg *Configuration) {
	a.initializeMiddlewares()
	a.initializeRoutes()
	a.initializeStatic()
}

// Run the listener
func (a *App) Run() {
	addr := a.Config.Host + ":" + a.Config.Port
	fmt.Println("Listen at: " + addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeMiddlewares() {
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.NoCache)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(middleware.RedirectSlashes)
}

func (a *App) initializeRoutes() {
	a.Router.Get("/", a.newPad)
	a.Router.Route("/{padname}", func(r chi.Router) {
		r.Get("/", a.getPad)
		r.Post("/", a.setPad)
	})
}

func (a *App) initializeStatic() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	static := filepath.Join(workDir, "static")

	a.Router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(static, "/favicon.ico"))
	})

	staticHandler(a.Router, "/static", http.Dir(static))
}
