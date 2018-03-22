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
	DB   string `default:"pad.db"`
	Salt string `default:""`
	Host string `default:"0.0.0.0"`
	Port string `default:"8080"`
}

var flagVersion = flag.Bool("version", false, "Show the version number and information")

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
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

	app, err := NewPadApp(&cfg)
	if err != nil {
		panic(err)
	}
	app.Initialize(&cfg)
	app.Run()
}
