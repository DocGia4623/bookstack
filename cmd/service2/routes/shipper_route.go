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
		shipperRouter.GET("", mw.AuthorizeRole(constant.ReadUser), shipperController.GetAllShipper)
		// Lấy tất cả order có địa chỉ trùng với nơi làm việc của shipper
		shipperRouter.GET("/orders", shipperController.GetOrderInRange)
		// Nhận đơn hàng
		shipperRouter.POST("/orders/:orderId", mw.AuthorizeRole(constant.ReceiveOrder), shipperController.ReceiveOrder)
		// Lấy tất cả đơn hàng đã  nhận
		shipperRouter.GET("/orders/received", mw.AuthorizeRole(constant.ReceiveOrder), shipperController.GetReceivedOrders)
	}
}
