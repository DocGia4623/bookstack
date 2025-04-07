package main

import (
	"bookstack/cmd/service2/routes"
	"bookstack/config"
	"bookstack/internal/wire"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config first
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Redis
	config.ConnectRedis(conf)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start listening for new orders
	app.ShipperController.StartListeningForNewOrders()

	routes.ShipperRoutes(router, app.Middleware, app.ShipperController)
	router.Run(":8081")
}
