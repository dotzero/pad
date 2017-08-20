package main

import (
	"html/template"
	"net/http"

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

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		cnt, err := app.GetNextCounter()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h := service.NewHash(app.Config.Salt)
		e, _ := h.EncodeInt64([]int64{cnt})

		http.Redirect(w, r, e, 301)
	})

	app.Router.Route("/{name}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			d := PadData{
				Name:    chi.URLParam(r, "name"),
				Content: chi.URLParam(r, "name"),
			}
			t := template.New("main")
			t, _ = t.ParseFiles("templates/main.html")
			t.Execute(w, d)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			context := r.Form.Get("context")
			w.Write([]byte(context))
		})
	})

	err := http.ListenAndServe(":"+cfg.Port, app.Router)
	if err != nil {
		panic(err)
	}
}
