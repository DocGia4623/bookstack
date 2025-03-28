package main

import (
	"bookstack/config"
	"bookstack/internal/repository"
	"bookstack/internal/wire"
	"bookstack/routes"
	"log"

	_ "bookstack/docs" // Import tài liệu Swagger đã được tạo

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // Thêm thư viện Swagger
)

// @title BookStack API
// @version 1.0
// @description This is a book ordering API built with Golang and Gin framework.
// @host localhost:8080
// @BasePath /
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

	// Seed database roles and permissions
	repository.SeedRolesAndPermissions()

	// Define routes
	routes.AuthRoute(*app.AuthenticationController, router)
	routes.UserRoute(*app.UserController, app.Middleware, router)
	routes.BookRoute(*app.BookController, router)
	routes.OrderRoute(*app.OrderController, router)

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	router.Run(":8080")
}
