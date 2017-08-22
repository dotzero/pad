package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

// Configuration is a Pad configuration
type Configuration struct {
	RedisURI    string `envconfig:"redis_uri" default:"redis://localhost:6379/0"`
	RedisPrefix string `envconfig:"redis_prefix" default:"pad"`
	Salt        string `default:"salt"`
	Port        string `default:"8080"`
}

func main() {
	var cfg Configuration
	if err := envconfig.Process("pad", &cfg); err != nil {
		panic(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app := &App{}
	app.Initialize(&cfg, workDir)
	app.Run()
}
