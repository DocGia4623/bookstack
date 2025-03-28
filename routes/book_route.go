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
		// tạo 1 cuốn sách, ko tạo chapter và page
		BookRoutes.POST("/", bookController.CreateBook)
		// Lấy tất cả sách
		BookRoutes.GET("/", bookController.GetBooks)
		// Thêm chapter cho 1 cuốn sách
		BookRoutes.POST("/:bookId/chapter", bookController.CreateChapter)
		// Lấy chapter của 1 cuốn sách
		BookRoutes.GET("/:bookId/chapter", bookController.GetChapters)
		// Thêm trang vào chapter
		BookRoutes.POST("/chapter/:chapterId/addpage", bookController.AddPage)
		// Lấy trang của chapter
		BookRoutes.POST("/chapter/:chapterId/getpage", bookController.GetPages)
		// Tạo kệ sách
		BookRoutes.POST("/shelve", bookController.CreateShelve)
	}
}
