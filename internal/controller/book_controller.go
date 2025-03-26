package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type BookController struct {
	bookSerivce service.BookService
	userService service.UserService
}

func NewBookController(service service.BookService, userServ service.UserService) *BookController {
	return &BookController{
		bookSerivce: service,
		userService: userServ,
	}
}

func (controller *BookController) CreateShelve(c *gin.Context) {
	var shelveRequest request.ShelveCreateRequest
	var webResponse response.WebResponse

	header := c.Request.Header.Get("Authorization")
	userId, err := controller.userService.GetUserIdByToken(header)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "no token found",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	user, err := controller.userService.GetUserById(userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant find user",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Kiểm tra dữ liệu đầu vào
	if err := c.ShouldBindJSON(&shelveRequest); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	shelve, err := controller.bookSerivce.CreateShelve(userId, shelveRequest)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "error during create shelve" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	var shelveResponse response.ShelveResponse
	copier.Copy(shelveResponse, shelve)
	shelveResponse.CreatedBy = user.FullName
	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Shelve created successfully",
		Data:    shelveResponse,
	}
	c.JSON(http.StatusCreated, webResponse)
}

func (controller *BookController) CreateBook(c *gin.Context) {
	var bookRequest request.BookCreateRequest
	var webResponse response.WebResponse

	header := c.Request.Header.Get("Authorization")
	userId, err := controller.userService.GetUserIdByToken(header)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "no token found",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	user, err := controller.userService.GetUserById(userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant find user",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
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

	book, err := controller.bookSerivce.CreateBook(userId, bookRequest)
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
	var bookReponse response.BookResponse
	copier.Copy(bookReponse, book)
	for _, tags := range book.Tags {
		bookReponse.Tags = append(bookReponse.Tags, tags.Name)
	}
	bookReponse.Shelve = book.Shelve.Name
	bookReponse.CreatedBy = user.FullName

	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Book created successfully",
		Data:    bookReponse,
	}
	c.JSON(http.StatusCreated, webResponse)
}
