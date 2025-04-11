package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func OrderRoute(controller controller.OrderController, router *gin.Engine) {
	OrderRoutes := router.Group("/order")
	{
		OrderRoutes.POST("/", controller.CreateOrder)
		OrderRoutes.POST("/paypal/:orderId", controller.CreatePaypalOrder)
		OrderRoutes.GET("/", controller.GetUserOrder)
		OrderRoutes.POST("/:orderId/cancel", controller.CancelOrder)
	}
	router.POST("/paypal/webhook", controller.HandlePaypalWebhook)
}
