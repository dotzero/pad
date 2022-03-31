package handlers

import (
	"bytes"
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
	h := hash.New("", 3)

	handler := Redirect(s, h)

	router := chi.NewRouter()
	router.Get("/", handler)

	w, err := testRequest(router, "GET", "/", "")
	is.NoErr(err)
	is.Equal(w.Code, http.StatusFound)
}

// func TestHandleGet(t *testing.T) {
// 	is := is.New(t)

// 	s := storage.NewMock()
// 	err := s.Set("foo", "bar")
// 	is.NoErr(err)

// 	handler := Get(s, "../web/templates/main.html")

// 	router := chi.NewRouter()
// 	router.Get("/{padname}", handler)

// 	w, err := testRequest(router, "GET", "/foo", "")
// 	is.NoErr(err)
// 	is.Equal(w.Code, http.StatusOK)
// 	is.True(strings.Contains(w.Body.String(), "bar"))
// }

func TestHandleRaw(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()
	err := s.Set("foo", "bar")
	is.NoErr(err)

	handler := Raw(s)

	router := chi.NewRouter()
	router.Get("/raw/{padname}", handler)

	w, err := testRequest(router, "GET", "/raw/foo", "")
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "bar"))
}

func TestHandleSet(t *testing.T) {
	is := is.New(t)

	s := storage.NewMock()

	handler := Set(s)

	router := chi.NewRouter()
	router.Post("/{padname}", handler)

	form := url.Values{}
	form.Set("t", "content")

	w, err := testRequest(router, "POST", "/foo", form.Encode())
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), `"message":"ok"`))
	is.True(strings.Contains(w.Body.String(), `"padname":"foo"`))
}

func testRequest(
	handler http.Handler,
	method string,
	address string,
	body string,
) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, address, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	return resp, nil
}
