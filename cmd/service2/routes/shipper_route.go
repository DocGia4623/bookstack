package routes

import (
	"bookstack/internal/constant"
	"bookstack/internal/controller"
	"bookstack/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ShipperRoutes(router *gin.Engine, mw *middleware.Middleware, shipperController *controller.ShipperController) {
	shipperRouter := router.Group("/shippers")
	{
		shipperRouter.GET("", mw.AuthorizeRole(constant.ReceiveOrder), shipperController.GetAllShipper)
		shipperRouter.GET("/orders", shipperController.GetOrderInRange)
	}
}
