package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dotzero/pad/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// App is a Pad app
type App struct {
	Config *Configuration
	Router *chi.Mux
	Redis  *service.Redis
	HashID *service.HashID
}

// Pad is a Pad data
type Pad struct {
	Name    string
	Content string
}

// Initialize is init the App
func (a *App) Initialize(cfg *Configuration, pwd string) {
	a.Config = cfg
	a.Router = chi.NewRouter()
	a.Redis = service.NewRedisClient(cfg.RedisURI, cfg.RedisPrefix)
	a.HashID = service.NewHashID(cfg.Salt, 3)
	a.initializeMiddlewares()
	a.initializeRoutes()
	a.initializeStatic(pwd)
}

// Run starts the listener
func (a *App) Run() {
	fmt.Printf("=> RedisURI: %s\n", a.Config.RedisURI)
	fmt.Printf("=> RedisPrefix: %s\n", a.Config.RedisPrefix)
	fmt.Println("Listen at: 0.0.0.0:" + a.Config.Port)
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, a.Router))
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

func (a *App) initializeStatic(pwd string) {
	static := filepath.Join(pwd, "static")

	a.Router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(static, "/favicon.ico"))
	})

	staticHandler(a.Router, "/static", http.Dir(static))
}

func (a *App) newPad(w http.ResponseWriter, r *http.Request) {
	cnt, err := a.Redis.GetNextCounter()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hash := a.HashID.Encode(cnt)
	http.Redirect(w, r, "/"+hash, 301)
}

func (a *App) getPad(w http.ResponseWriter, r *http.Request) {
	padname := chi.URLParam(r, "padname")
	content := a.Redis.GetPad(padname)

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

	err = a.Redis.SetPad(padname, content)
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

func staticHandler(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
