package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func BookRoute(bookController controller.BookController, router *gin.Engine) {
	BookRoutes := router.Group("/book")
	{
		// 1 cuốn sách hoàn chỉnh
		BookRoutes.POST("/complete", bookController.CreateCompleteBook)
		BookRoutes.POST("/", bookController.CreateBook)
		BookRoutes.GET("/", bookController.GetBooks)
		BookRoutes.POST("/:bookId/chapter", bookController.CreateChapter)
		BookRoutes.GET("/:bookId/chapter", bookController.GetChapters)
		BookRoutes.POST("/chapter/:chapterId/addpage", bookController.AddPage)
		BookRoutes.POST("/chapter/:chapterId/getpage", bookController.GetPages)
		BookRoutes.POST("/shelve", bookController.CreateShelve)
	}
}
