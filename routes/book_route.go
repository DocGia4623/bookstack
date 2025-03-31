package routes

import (
	"bookstack/internal/controller"

	"github.com/gin-gonic/gin"
)

func BookRoute(bookController controller.BookController, router *gin.Engine) {
	BookRoutes := router.Group("/book")
	{
		//book
		BookRoutes.POST("/complete", bookController.CreateCompleteBook)
		BookRoutes.POST("/", bookController.CreateBook)
		BookRoutes.GET("/", bookController.GetBooks)
		BookRoutes.PUT("/:bookId", bookController.UpdateBook)
		BookRoutes.DELETE("/:bookId", bookController.DeleteBook)
		//shelve
		BookRoutes.POST("/shelve", bookController.CreateShelve)
		BookRoutes.GET("/shelve", bookController.GetShelves)
		BookRoutes.DELETE("/shelve/:shelveId", bookController.DeleteShelve)
		//chapter
		BookRoutes.POST("/:bookId/chapter", bookController.CreateChapter)
		BookRoutes.GET("/:bookId/chapter", bookController.GetChapters)
		BookRoutes.PUT("/:bookId/chapter/:chapterId", bookController.UpdateChapter)
		BookRoutes.DELETE("/:bookId/chapter/:chapterId", bookController.DeleteChapter)
		//page
		BookRoutes.POST("/chapter/:chapterId/page", bookController.AddPage)
		BookRoutes.GET("/chapter/:chapterId/page", bookController.GetPages)
		BookRoutes.PUT("/chapter/:chapterId/page/:pageId", bookController.UpdatePage)
		BookRoutes.DELETE("/chapter/:chapterId/page/:pageId", bookController.DeletePage)
	}
}
