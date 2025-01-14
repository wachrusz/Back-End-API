package main

import (
	"log"
	"os"

	"github.com/wachrusz/Back-End-API/internal/app"
	"github.com/wachrusz/Back-End-API/internal/config"
)

// @title			Cash Advisor API
// @version         0.1
// @description     Backend API for managing user profiles, authentication, analytics, and more.
// @termsOfService  http://swagger.io/terms/

// @contact.name   	CADV Support
// @contact.name	Mikhail Vakhrushin
// @contact.email 	wachrusz@gmail.com

// @host      	localhost:8080
// @schemes 	https
// @BasePath  	/v1

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description To authorize,
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
