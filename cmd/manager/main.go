package main

import (
	"math/rand"
	"time"

	"github.com/netmoth/netmoth/internal/web"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	web.New()
}
