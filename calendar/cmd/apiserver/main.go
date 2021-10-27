package main

import (
	"calendar/internal/app/httpserver"
	"calendar/internal/config"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"runtime"
)

var (
	configPath string
)

func init() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := config.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := httpserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
