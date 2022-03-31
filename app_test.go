package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/zero-pkg/tpl"
)

func TestNew(t *testing.T) {
	is := is.New(t)

	tmpdir, err := os.MkdirTemp("", "padtest")
	is.NoErr(err)

	defer os.RemoveAll(tmpdir)

	app, err := New(Opts{
		DatabasePath: tmpdir,
		AssetsPath:   "web",
	})
	is.NoErr(err)
	is.Equal(app.Storage.Path(), filepath.Join(tmpdir, databaseFile))
	is.True(app.HashEncoder != nil)
	is.True(app.Templates != nil)
}

func TestRoutes(t *testing.T) {
	is := is.New(t)

	app := &App{
		Templates: tpl.New(),
	}

	routes := app.routes()

	is.Equal(len(routes.Middlewares()), 5)
	is.Equal(len(routes.Routes()), 7)
}
