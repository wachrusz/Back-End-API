package main

import (
	"github.com/wachrusz/Back-End-API/internal/app"
	"github.com/wachrusz/Back-End-API/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}
	if err = app.Run(cfg); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
	os.Exit(0)
}
