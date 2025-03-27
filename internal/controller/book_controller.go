package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/models"
	"bookstack/internal/service"
	"net/http"
	"strconv"

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

func (controller *BookController) GetBooks(c *gin.Context) {
	var webResponse response.WebResponse
	books, err := controller.bookSerivce.GetAllBook()
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "error during get book" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	var booksResponse []response.BookResponse
	// Copy dữ liệu từng cuốn sách vào response
	for _, book := range books {
		var bookResponse response.BookResponse
		copier.Copy(&bookResponse, &book)
		booksResponse = append(booksResponse, bookResponse)
	}
	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Books",
		Data:    booksResponse,
	}
	c.JSON(http.StatusCreated, webResponse)
}

func (controller *BookController) CreateChapter(c *gin.Context) {
	var webResponse response.WebResponse
	var request request.BookChapterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	chapter, err := controller.bookSerivce.CreateChapter(request)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Cant create" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Created chapter",
		Data:    chapter,
	}
	c.JSON(http.StatusCreated, webResponse)
}

func (controller *BookController) GetChapters(c *gin.Context) {
	var webResponse response.WebResponse

	// Lấy bookID từ form-data
	bookIDStr := c.PostForm("bookId") // Hoặc c.DefaultPostForm("bookId", "0")
	if bookIDStr == "" {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "missing bookId in form-data",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Chuyển bookID từ string -> int
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "invalid bookId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Gọi service để lấy chapters của book
	chapters, err := controller.bookSerivce.GetChaptersOfBook(bookID)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "error getting chapters: " + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	// Trả về danh sách chương
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Chapters retrieved successfully",
		Data:    chapters,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *BookController) AddPage(c *gin.Context) {
	var webResponse response.WebResponse
	var request request.PageRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "request invalid",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	page, err := controller.bookSerivce.AddPage(request)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Cant add page",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	// Trả về danh sách chương
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Page created",
		Data:    page,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *BookController) GetPages(c *gin.Context) {
	var webResponse response.WebResponse
	var pages []models.Page
	// Lấy bookID từ form-data
	ChapterIdstr := c.PostForm("pageId") // Hoặc c.DefaultPostForm("bookId", "0")
	if ChapterIdstr == "" {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "missing ChapterId in form-data",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	ChapterId, err := strconv.Atoi(ChapterIdstr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Chapterid must be integer",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	pages, err = controller.bookSerivce.GetPageChapter(ChapterId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Error from server",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	// Trả về danh sách chương
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pages",
		Data:    pages,
	}
	c.JSON(http.StatusOK, webResponse)
}
