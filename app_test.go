package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/matryer/is"

	"github.com/dotzero/pad/hash"
	"github.com/dotzero/pad/storage"
)

func TestHandleRedirect(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()
	router := mockRouter(s)

	r, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusFound)
}

func TestHandleGet(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()
	err := s.Set("foo", "bar")
	is.NoErr(err)

	router := mockRouter(s)

	r, err := http.NewRequest("GET", "/foo", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "bar"))
}

func TestHandleRaw(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()
	err := s.Set("foo", "bar")
	is.NoErr(err)

	router := mockRouter(s)

	r, err := http.NewRequest("GET", "/raw/foo", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "bar"))
}

func TestHandleSet(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()
	router := mockRouter(s)

	form := url.Values{}
	form.Set("t", "content")

	r, err := http.NewRequest("POST", "/foo", strings.NewReader(form.Encode()))
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), `"message":"ok"`))
	is.True(strings.Contains(w.Body.String(), `"padname":"foo"`))
}

func mockRouter(s *storage.MockStorage) chi.Router {
	app := App{
		Storage:     s,
		HashEncoder: hash.New("", 3),
		Opts: Opts{
			AssetsPath: "./web",
		},
	}

	return app.routes()
}
