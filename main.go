package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"

	"github.com/dotzero/pad/app"
)

var (
	// These variables should be set by the linker when compiling
	// CommitHash is the git hash of last commit
	CommitHash = "Unknown"
	// CompileDate is the date of build
	CompileDate = "Unknown"
)

// Opts with command line flags and env
type Opts struct {
	Host       string `long:"host" env:"PAD_HOST" default:"0.0.0.0" description:"listening address"`
	Port       int    `long:"port" env:"PAD_PORT" default:"8080" description:"listening port"`
	BoltPath   string `long:"bolt-path" env:"BOLT_PATH" default:"./var" description:"BoltDB data directory"`
	SecretKey  string `long:"secret" env:"PAD_SECRET" description:"shared secret key used to generate IDs"`
	StaticPath string `long:"static-path" env:"STATIC_PATH" default:"./static" description:"path to static assets"`
	TmlPath    string `long:"tpl-path" env:"TPL_PATH" default:"./templates" description:"path to template files"`
	TplExt     string `long:"tpl-ext" env:"TPL_EXT" default:".html" description:"template file extension"`
	Verbose    bool   `long:"verbose" description:"enable verbose logging"`
	Version    bool   `short:"v" long:"version" description:"show version information"`
}

func main() {
	var opts Opts

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("Commit hash: %s\nCompile date: %s\n", CommitHash, CompileDate)
		os.Exit(0)
	}

	setupLog(opts.Verbose)
	log.Printf("[DEBUG] opts: %+v", opts)

	app, err := app.New(app.CommonOpts{
		BoltPath:   opts.BoltPath,
		SecretKey:  opts.SecretKey,
		StaticPath: opts.StaticPath,
		TmlPath:    opts.TmlPath,
		TplExt:     opts.TplExt,
	})
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	if err := app.Run(context.Background(), opts.Host, opts.Port); err != nil {
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
