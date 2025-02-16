package main

import (
	"log"

	"github.com/Avazbek-02/Online-Hotel-System/config"
	"github.com/Avazbek-02/Online-Hotel-System/internal/app"
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
