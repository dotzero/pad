package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotzero/pad/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
)

// Configuration is a Pad configuration
type Configuration struct {
	RedisURI    string `envconfig:"redis_uri" default:"redis://localhost:6379/0"`
	RedisPrefix string `envconfig:"redis_prefix" default:"pad"`
	Salt        string `default:"salt"`
	Port        string `default:"8080"`
}

// PadApp is a Pad app
type PadApp struct {
	Config *Configuration
	Router *chi.Mux
	Redis  *redis.Client
}

// PadData is a Pad data
type PadData struct {
	Name    string
	Content string
}

type Response struct {
	Message string `json:"message"`
	PadName string `json:"padname"`
}

func (c *PadApp) prefixed(key string) string {
	return c.Config.RedisPrefix + ":" + key
}

// GetNextCounter returns uniq hash
func (c *PadApp) GetNextCounter() (int64, error) {
	val, err := c.Redis.Incr(c.prefixed("counter")).Result()
	if err != nil {
		return 0, err
	}

	return val, nil
}

var app *PadApp

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
		panic(err)
	}

	app = &PadApp{
		Config: &cfg,
		Router: chi.NewRouter(),
		Redis:  service.NewRedisClient(cfg.RedisURI),
	}

	app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.NoCache)
	app.Router.Use(middleware.RealIP)
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.RedirectSlashes)

	app.Router.Get("/", routeMain)
	app.Router.Route("/{name}", func(r chi.Router) {
		r.Get("/", routeLoad)
		r.Post("/", routeSave)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	routeStatic(app.Router, "/static", http.Dir(filesDir))

	log.Printf("Redis uri: %s\n", cfg.RedisURI)
	log.Printf("Redis prefix: %s\n", cfg.RedisPrefix)
	log.Printf("Listen at: http://0.0.0.0:%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, app.Router); err != nil {
		panic(err)
	}
}

func routeMain(w http.ResponseWriter, r *http.Request) {
	cnt, err := app.GetNextCounter()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h := service.NewHash(app.Config.Salt)
	e, _ := h.EncodeInt64([]int64{cnt})

	http.Redirect(w, r, e, 301)
}

func routeLoad(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value, err := app.Redis.Get(name).Result()
	if err != nil {
		value = ""
	}
	d := PadData{
		Name:    name,
		Content: value,
	}

	t := template.New("main")
	t, _ = t.ParseFiles("templates/main.html")
	t.Execute(w, d)
}

func routeSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := chi.URLParam(r, "name")
	context := r.Form.Get("t")
	_ = app.Redis.Set(name, context, 0).Err()
	res, _ := json.Marshal(Response{
		Message: "ok",
		PadName: name,
	})
	w.Write(res)
}

func routeStatic(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

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
