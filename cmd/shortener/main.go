package main

import (
	"github.com/AntonNikol/go-shortener/internal/app/server"
	"github.com/AntonNikol/go-shortener/internal/config"
)

func main() {
	cfg := config.Get()
	server.Run(cfg)
}
