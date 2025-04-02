package routes

import (
	"bookstack/cmd/service2/internal/controller.go"

	"github.com/gin-gonic/gin"
)

func ShipperRoutes(router *gin.Engine, shipperController *controller.ShipperController) {
	shipperRouter := router.Group("/shippers")
	{
		shipperRouter.GET("", shipperController.GetAllShipper)
		shipperRouter.GET("/orders", shipperController.GetOrderInRange)
	}
}
