package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/zero-pkg/tpl"

	"github.com/dotzero/pad/app/hash"
	"github.com/dotzero/pad/app/storage"
)

const (
	boltFile = "pad.db"
)

// App is a Pad application
type App struct {
	CommonOpts
	Storage     *storage.BoltStorage
	HashEncoder *hash.Hash
	Templates   *tpl.Templates
	httpServer  *http.Server
}

// New prepares application
func New(commonOpts CommonOpts) (*App, error) {
	app := &App{}
	app.SetCommon(commonOpts)

	if err := app.setupDataStore(); err != nil {
		return nil, err
	}

	if err := app.setupTpl(); err != nil {
		return nil, err
	}

	app.HashEncoder = hash.New(commonOpts.SecretKey, 3)

	return app, nil
}

// Run listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (a *App) Run(ctx context.Context, address string, port int) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.httpServer = a.makeHTTPServer(address, port, a.routes())

	go func() {
		addr := fmt.Sprintf("%s:%d", address, port)
		log.Printf("[INFO] http server listen at: http://" + addr)

		err := a.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] failed to listen and serve, %+v", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(timeoutCtx)
}

// setupDataStore initializes BoltDB store
func (a *App) setupDataStore() error {
	if err := makeDirs(a.BoltPath); err != nil {
		return err
	}

	store, err := storage.New(filepath.Join(a.BoltPath, boltFile))
	if err != nil {
		return err
	}

	a.Storage = store

	return nil
}

// setupTpl initializes Templates engine
func (a *App) setupTpl() error {
	templ, err := tpl.New().ParseDir(a.TmlPath, a.TplExt)
	if err != nil {
		return err
	}

	a.Templates = templ

	return nil
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
			if e := os.MkdirAll(dir, 0700); e != nil { // nolint
				return fmt.Errorf("can't make directory %s", dir)
			}
		}
	}

	return nil
}
