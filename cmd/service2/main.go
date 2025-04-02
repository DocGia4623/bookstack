package main

import (
	"bookstack/cmd/service2/routes"
	"bookstack/cmd/service2/wire"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	routes.ShipperRoutes(router, app.ShipperController)
	router.Run(":8081")
}
