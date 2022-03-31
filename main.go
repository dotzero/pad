package main

import (
	"fmt"
	"log"
	"os"

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
	Host         string `long:"host" env:"PAD_HOST" default:"0.0.0.0" description:"listening address"`
	Port         int    `long:"port" env:"PAD_PORT" default:"8080" description:"listening port"`
	DatabasePath string `long:"db" env:"PAD_DB_PATH" default:"db" description:"path to database files"`
	AssetsPath   string `long:"path" env:"PAD_ASSETS_PATH" default:"web" description:"path to web assets"`
	SecretKey    string `long:"secret" env:"PAD_SECRET" description:"the shared secret key used to generate ids"`
	Verbose      bool   `long:"verbose" description:"verbose logging"`
	Version      bool   `short:"v" long:"version" description:"show the version number"`
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("Version: %s\nCommit hash: %s\nCompile date: %s\n", Version, CommitHash, CompileDate)
		os.Exit(0)
	}

	setupLog(opts.Verbose)
	log.Printf("[DEBUG] opts: %+v", opts)

	app, err := New(opts)
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("[WARN] http server terminated, %s", err)
	}
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
