package main

import (
	"flag"
	"log"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/sensor"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "cfg", "cmd/agent/config.yml", "configuration file path")
	flag.Parse()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	sensor.New(cfg)
}
