package main

import (
	"bookstack/config"
	"bookstack/internal/wire"
	"bookstack/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	conf, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	// Connect to the database
	config.Connect(conf)
	// Initialize the router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Set up routes
	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	routes.AuthRoute(*app.AuthenticationController, router)
	routes.UserRoute(*app.UserController, router)
	router.Run(":8080")
}
