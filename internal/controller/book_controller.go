package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookSerivce service.BookService
}

func NewBookController(service service.BookService) *BookController {
	return &BookController{
		bookSerivce: service,
	}
}

func (controller *BookController) CreateBook(c *gin.Context) {
	var bookRequest request.BookCreateRequest
	var webResponse response.WebResponse

	// Kiểm tra dữ liệu đầu vào
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	book, err := controller.bookSerivce.CreateBook(bookRequest)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "error during create book" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Book created successfully",
		Data:    book,
	}
	c.JSON(http.StatusCreated, webResponse)
}
