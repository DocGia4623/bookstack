package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoute(authController controller.AuthenticationController, router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/logout", authController.Logout)
	}
}
