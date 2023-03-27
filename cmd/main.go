package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/sensor"
	"github.com/netmoth/netmoth/internal/web"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var configPath string
	flag.StringVar(&configPath, "cfg", "config.yml", "configuration file path")
	flag.Parse()

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	go web.New()
	sensor.New(ctx, cfg)
}
