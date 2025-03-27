package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func BookRoute(bookController controller.BookController, router *gin.Engine) {
	BookRoutes := router.Group("/book")
	{
		BookRoutes.POST("/create", bookController.CreateBook)
		BookRoutes.POST("/createshelve", bookController.CreateShelve)
		BookRoutes.GET("/", bookController.GetBooks)
		BookRoutes.POST("/createchapter", bookController.CreateChapter)
		BookRoutes.POST("/getchapter", bookController.GetChapters)
		BookRoutes.POST("/addpage", bookController.AddPage)
		BookRoutes.POST("/getpage", bookController.GetPages)
	}
}
