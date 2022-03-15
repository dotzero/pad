package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dotzero/pad/hash"
	"github.com/dotzero/pad/storage"
	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"
)

var (
	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version = "0.0.0-unknown"
	// CommitHash is the git hash of last commit
	CommitHash = "Unknown"
	// CompileDate is the date of build
	CompileDate = "Unknown"
)

// Opts with command line flags and env
type Opts struct {
	BoltPath  string `long:"db" env:"PAD_DB_PATH" default:"./db" description:"path to database"`
	SecretKey string `long:"secret" env:"PAD_SECRET" description:"secret key"`

	Host    string `long:"host" env:"PAD_HOST" default:"0.0.0.0" description:"host"`
	Port    int    `long:"port" env:"PAD_PORT" default:"8080" description:"port"`
	WebPath string `long:"path" env:"PAD_PATH" default:"./web" description:"path to web assets"`

	Verbose bool `short:"v" long:"verbose" description:"enable verbose logging"`
	Version bool `long:"version" description:"show the version number and information"`
}

type hashEncoder interface {
	Encode(num int64) string
}

type padStorage interface {
	Get(name string) (value string, err error)
	Set(name string, value string) error
	NextCounter() (next uint64, err error)
}

// App is a Pad app
type App struct {
	Opts
	Storage     padStorage
	HashEncoder hashEncoder
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	setupLog(opts.Verbose)
	log.Printf("[DEBUG] opts: %+v", opts)

	if opts.Version {
		// If -version was passed
		fmt.Printf("Version: %s\nCommit hash: %s\nCompile date: %s\n", Version, CommitHash, CompileDate)
		os.Exit(0)
	}

	app, err := New(opts)
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	err = app.Run()
	log.Fatalf("[WARN] http server terminated, %s", err)
}

// New prepares application and return it
func New(opts Opts) (*App, error) {
	if err := makeDirs(opts.BoltPath); err != nil {
		return nil, err
	}

	boltBackend, err := storage.New(opts.BoltPath, "pad.db")
	if err != nil {
		return nil, err
	}

	return &App{
		Opts:        opts,
		Storage:     boltBackend,
		HashEncoder: hash.New(opts.SecretKey, 3),
	}, nil
}

// Run the listener
func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.Opts.Host, a.Opts.Port)
	log.Printf("[INFO] http server listen at: http://" + addr)

	router := a.routes()
	return http.ListenAndServe(addr, router)
}

func setupLog(verbose bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	log.SetOutput(filter)
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
			if e := os.MkdirAll(dir, 0700); e != nil {
				return fmt.Errorf("can't make directory %s", dir)
			}
		}
	}
	return nil
}
