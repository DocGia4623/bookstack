package routes

import (
	"bookstack/internal/constant"
	"bookstack/internal/controller"
	"bookstack/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(controller controller.UserController, mw *middleware.Middleware, router *gin.Engine) {
	UserRoutes := router.Group("/user")
	{
		// Lấy tất cả user
		UserRoutes.GET("/", mw.AuthorizeRole(constant.ReadUser), controller.GetAllUser)
		// Update user
		UserRoutes.PUT("/", controller.UpdateUser)
	}
}
