package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func UserRoute(controller controller.UserController, router *gin.Engine) {
	UserRoutes := router.Group("/user")
	{
		UserRoutes.GET("/", controller.GetAllUser)
		UserRoutes.PUT("/", controller.UpdateUser)
	}
}
