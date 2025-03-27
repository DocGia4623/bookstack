package main

import (
	"bookstack/config"
	"bookstack/internal/repository"
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
	repository.SeedRolesAndPermissions()
	routes.AuthRoute(*app.AuthenticationController, router)
	routes.UserRoute(*app.UserController, app.Middleware, router)
	routes.BookRoute(*app.BookController, router)
	routes.OrderRoute(*app.OrderController, router)
	router.Run(":8080")
}
