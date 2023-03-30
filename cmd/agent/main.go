package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/sensor"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var configPath string
	flag.StringVar(&configPath, "cfg", "config.yml", "configuration file path")
	flag.Parse()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	sensor.New(cfg)
}
