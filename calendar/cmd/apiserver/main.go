package main

import (
	"calendar/internal/app/httpserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/httpserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := httpserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := httpserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
