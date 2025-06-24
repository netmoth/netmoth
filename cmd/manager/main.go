package main

import (
	"github.com/netmoth/netmoth/internal/web"
)

func main() {
	// Start the web server with local config
	web.New("cmd/manager/config.yml")
}
