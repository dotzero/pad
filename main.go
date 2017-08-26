package main

import (
	"flag"
	"fmt"
	"os"

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
	RedisURI    string `envconfig:"redis_uri" default:"redis://localhost:6379/0"`
	RedisPrefix string `envconfig:"redis_prefix" default:"pad"`
	Salt        string `default:"salt"`
	Port        string `default:"8080"`
}

var flagVersion = flag.Bool("version", false, "Show the version number and information")

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
		panic(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.Parse()
	if *flagVersion {
		// If -version was passed
		fmt.Println("Version:", Version)
		fmt.Println("Commit hash:", CommitHash)
		fmt.Println("Compile date", CompileDate)
		os.Exit(0)
	}

	app := &App{}
	app.Initialize(&cfg, workDir)
	app.Run()
}
