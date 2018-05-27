package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dotzero/pad/service"
	"github.com/go-chi/chi"

	"github.com/dotzero/pad/mocks"
	"github.com/matryer/is"
)

func TestHandleNewPad(t *testing.T) {
	is := is.New(t)
	router := NewTestRouter()

	r, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusFound)
}

func TestHandleGetPad(t *testing.T) {
	is := is.New(t)
	router := NewTestRouter()

	r, err := http.NewRequest("GET", "/foobar", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "foobar"))
}

func TestHandleSetPad(t *testing.T) {
	is := is.New(t)
	router := NewTestRouter()

	form := url.Values{}
	form.Set("t", "content")

	r, err := http.NewRequest("POST", "/foobar", strings.NewReader(form.Encode()))
	is.NoErr(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), `"message":"ok"`))
	is.True(strings.Contains(w.Body.String(), `"padname":"foobar"`))
}

func NewTestRouter() chi.Router {
	app := App{
		Storage: mocks.NewStorage(),
		Unique:  service.NewHashID("", 3),
	}

	return app.routes()
}
