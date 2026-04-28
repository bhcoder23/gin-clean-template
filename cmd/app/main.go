package main

import (
	"log"

	"github.com/bhcoder23/gin-clean-template/config"
	"github.com/bhcoder23/gin-clean-template/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
