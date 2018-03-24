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

var flagSet = flag.NewFlagSet("", flag.ExitOnError)
var flagSilent = flagSet.Bool("silent", false, "Operate without emitting any output")
var flagVersion = flagSet.Bool("version", false, "Show the version number and information")

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
		panic(err)
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
		fmt.Println("Configuration")
		fmt.Println("=> DB:", cfg.DB)
		fmt.Println("=> Salt:", cfg.Salt)
		fmt.Println("=> Host:", cfg.Host)
		fmt.Println("=> Port:", cfg.Port)
	}

	app, err := NewPadApp(&cfg)
	if err != nil {
		panic(err)
	}
	app.Initialize(&cfg)
	app.Run(*flagSilent)
}
