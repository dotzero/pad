package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dotzero/pad/service"
	"github.com/kelseyhightower/envconfig"
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

// Configuration is a Pad configuration
type Configuration struct {
	DB      string `default:"pad.db"`
	Salt    string `default:""`
	Host    string `default:"0.0.0.0"`
	Port    string `default:"8080"`
	WebRoot string `default:"."`
}

// App is a Pad app
type App struct {
	Config      *Configuration
	BoltBackend *service.BoltBackend
	HashID      *service.HashID
}

var flagSet = flag.NewFlagSet("", flag.ExitOnError)
var flagSilent = flagSet.Bool("silent", false, "Operate without emitting any output")
var flagVersion = flagSet.Bool("version", false, "Show the version number and information")

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
		os.Exit(1)
	}

	flagSet.Parse(os.Args[1:])
	if *flagVersion {
		// If -version was passed
		fmt.Println("Version:", Version)
		fmt.Println("Commit hash:", CommitHash)
		fmt.Println("Compile date", CompileDate)
		os.Exit(0)
	}
	if *flagSilent == false {
		// If -silent was not passed
		log.Printf("Env DB: %s", cfg.DB)
		log.Printf("Env Salt: %s", cfg.Salt)
		log.Printf("Env Host: %s", cfg.Host)
		log.Printf("Env Port: %s", cfg.Port)
	}

	app, err := New(&cfg)
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}
	app.Run(*flagSilent)
}

// New prepares application and return it
func New(cfg *Configuration) (*App, error) {
	boltBackend, err := service.NewBoltBackend(cfg.DB)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:      cfg,
		BoltBackend: boltBackend,
		HashID:      service.NewHashID(cfg.Salt, 3),
	}, nil
}

// Run the listener
func (a *App) Run(flagSilent bool) {
	addr := a.Config.Host + ":" + a.Config.Port
	if flagSilent == false {
		fmt.Println("Listen at: http://" + addr)
	}
	router := a.routes()
	log.Fatal(http.ListenAndServe(addr, router))
}
