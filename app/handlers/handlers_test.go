package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/matryer/is"
)

func TestHandleRedirect(t *testing.T) {
	is := is.New(t)
	s := &storageMock{
		NextCounterFunc: func() (uint64, error) {
			return 1, nil
		},
	}
	e := &encoderMock{
		EncodeFunc: func(num int64) string {
			return "foo"
		},
	}

	handler := Redirect(s, e)

	router := chi.NewRouter()
	router.Get("/", handler)

	w, err := testRequest(router, "GET", "/", "")
	is.NoErr(err)
	is.Equal(w.Code, http.StatusFound)
}

func TestHandleGet(t *testing.T) {
	is := is.New(t)
	s := &storageMock{
		GetFunc: func(name string) (string, error) {
			return "bar", nil
		},
	}

	tpl := template.Must(template.New("").Parse("{{ .Padname }}{{ .Content }}"))

	handler := Get(s, tpl)

	router := chi.NewRouter()
	router.Get("/{padname}", handler)

	w, err := testRequest(router, "GET", "/foo", "")
	is.NoErr(err)
	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "foobar"))
}

func TestHandleRaw(t *testing.T) {
	is := is.New(t)
	s := &storageMock{
		GetFunc: func(name string) (string, error) {
			return "bar", nil
		},
	}

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
	s := &storageMock{
		SetFunc: func(name string, value string) error {
			return nil
		},
	}

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
