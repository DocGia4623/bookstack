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

// CreateShelve godoc
// @Summary Create a new shelve
// @Description Create a new shelve for a user
// @Tags Shelve
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param shelve body request.ShelveCreateRequest true "Shelve request body"
// @Success 201 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /shelves [post]
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
	copier.Copy(&shelveResponse, shelve)
	shelveResponse.CreatedBy = user.FullName
	var tags []response.TagResponse
	var tagResponsen response.TagResponse
	for _, tag := range shelve.Tags {
		err := copier.Copy(&tagResponsen, tag)
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
		tags = append(tags, tagResponsen)
	}
	shelveResponse.Tags = tags
	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Shelve created successfully",
		Data:    shelveResponse,
	}
	c.JSON(http.StatusCreated, webResponse)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book associated with a user
// @Tags Book
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param book body request.BookCreateRequest true "Book request body"
// @Success 201 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books [post]
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
	var tags []response.TagResponse
	var tag response.TagResponse
	copier.Copy(&bookReponse, book)
	for _, t := range book.Tags {
		copier.Copy(&tag, t)
		tags = append(tags, tag)
	}
	bookReponse.Shelve = book.Shelve.Name
	bookReponse.CreatedBy = user.FullName
	bookReponse.Tags = tags
	// Phản hồi thành công
	webResponse = response.WebResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Book created successfully",
		Data:    bookReponse,
	}
	c.JSON(http.StatusCreated, webResponse)
}

// GetBooks godoc
// @Summary Get all books
// @Description Retrieve all books available in the system
// @Tags Book
// @Produce json
// @Success 200 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books [get]
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
		// Kiểm tra CreatedBy có tồn tại không
		user, err := controller.userService.GetUserById(int(book.CreatedBy))
		if err != nil {
			webResponse = response.WebResponse{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "cant find user",
				Data:    nil,
			}
			c.JSON(http.StatusInternalServerError, webResponse)
			return
		}
		bookResponse.CreatedBy = user.FullName
		bookResponse.Shelve = book.Shelve.Name
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

// CreateChapter godoc
// @Summary Create a chapter for a book
// @Description Add a new chapter to an existing book
// @Tags Chapter
// @Accept json
// @Produce json
// @Param bookId path int true "Book ID"
// @Param chapter body request.BookChapterRequest true "Chapter request body"
// @Success 201 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books/{bookId}/chapters [post]
func (controller *BookController) CreateChapter(c *gin.Context) {
	bookIdStr := c.Param("bookId")                        // Lấy bookId từ URL
	bookId64, err := strconv.ParseUint(bookIdStr, 10, 32) // Chuyển thành uint64
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid bookId"})
		return
	}
	bookId := uint(bookId64) // Ép kiểu thành uint
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
	chapter, err := controller.bookSerivce.CreateChapter(bookId, request)
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

// GetChapters godoc
// @Summary Get all chapters of a book
// @Description Retrieve all chapters associated with a book
// @Tags Chapter
// @Produce json
// @Param bookId path int true "Book ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books/{bookId}/chapters [get]
func (controller *BookController) GetChapters(c *gin.Context) {
	var webResponse response.WebResponse
	bookIdStr := c.Param("bookId") // Lấy bookId từ URL
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get bookId" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Gọi service để lấy chapters của book
	chapters, err := controller.bookSerivce.GetChaptersOfBook(bookId)
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

// AddPage godoc
// @Summary Add a page to a chapter
// @Description Create a new page inside a given chapter
// @Tags Page
// @Accept json
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Param page body request.PageRequest true "Page request body"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /chapters/{chapterId}/pages [post]
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
	chapterIdStr := c.Param("chapterId")                        // Lấy bookId từ URL
	chapterId64, err := strconv.ParseUint(chapterIdStr, 10, 32) // Chuyển thành uint64
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid bookId"})
		return
	}
	chapterId := uint(chapterId64) // Ép kiểu thành uint

	page, err := controller.bookSerivce.AddPage(chapterId, request)
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

// GetPages godoc
// @Summary Get pages of a chapter
// @Description Retrieve all pages within a chapter
// @Tags Page
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /chapters/{chapterId}/pages [get]
func (controller *BookController) GetPages(c *gin.Context) {
	var webResponse response.WebResponse
	var pages []models.Page
	chapterIdStr := c.Param("chapterId")

	ChapterId, err := strconv.Atoi(chapterIdStr)
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
	// Trả về danh sách Trang
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pages",
		Data:    pages,
	}
	c.JSON(http.StatusOK, webResponse)
}

// CreateCompleteBook godoc
// @Summary Create a complete book with chapters and pages
// @Description Create a fully detailed book with associated chapters and pages
// @Tags Book
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param book body request.CompleteBookCreateRequest true "Complete book request body"
// @Success 201 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books/complete [post]
func (controller *BookController) CreateCompleteBook(c *gin.Context) {
	var webResponse response.WebResponse
	var request request.CompleteBookCreateRequest

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
	if err := c.ShouldBindJSON(&request); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	book, err := controller.bookSerivce.CreateCompleteBook(user.ID, request)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Server error",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	var bookResponse response.BookResponse
	copier.Copy(&bookResponse, book)
	bookResponse.CreatedBy = user.FullName
	bookResponse.Shelve = book.Shelve.Name
	var tags []response.TagResponse
	var tag response.TagResponse
	for _, t := range book.Tags {
		copier.Copy(&tag, t)
		tags = append(tags, tag)
	}
	bookResponse.Tags = tags
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pages",
		Data:    bookResponse,
	}
	c.JSON(http.StatusOK, webResponse)
}

// GetShelves godoc
// @Summary Get all shelves
// @Description Retrieve all shelves available in the system
// @Tags Shelve
// @Produce json
// @Success 200 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /shelves [get]
func (controller *BookController) GetShelves(c *gin.Context) {
	var webResponse response.WebResponse
	shelves, err := controller.bookSerivce.GetShelves()
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Server error",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var shelveResponses []response.ShelveResponse
	err = copier.Copy(&shelveResponses, shelves)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Server error" + err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	for _, shelveResponse := range shelveResponses {
		var tags []response.TagResponse
		var tagResponse response.TagResponse

		for _, tag := range shelveResponse.Tags {
			err := copier.Copy(&tagResponse, tag)
			if err != nil {
				webResponse = response.WebResponse{
					Code:    http.StatusInternalServerError,
					Status:  "error",
					Message: "Server error" + err.Error(),
					Data:    nil,
				}
				c.JSON(http.StatusInternalServerError, webResponse)
				return
			}
			tags = append(tags, tagResponse)
		}

		shelveResponse.Tags = tags
	}

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Pages",
		Data:    shelveResponses,
	}
	c.JSON(http.StatusOK, webResponse)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID
// @Tags Book
// @Produce json
// @Param bookId path int true "Book ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books/{bookId} [delete]
func (controller *BookController) DeleteBook(c *gin.Context) {
	var webResponse response.WebResponse
	bookIdStr := c.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get bookId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.bookSerivce.DeleteBook(bookId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant delete book",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "book deleted",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book by ID
// @Tags Book
// @Accept json
// @Produce json
// @Param bookId path int true "Book ID"
// @Param book body request.BookCreateRequest true "Book request body"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /books/{bookId} [put]
func (controller *BookController) UpdateBook(c *gin.Context) {
	var webResponse response.WebResponse
	var request request.BookCreateRequest
	bookIdStr := c.Param("bookId")
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get bookId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	book, err := controller.bookSerivce.UpdateBook(bookId, request)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant update book",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var bookResponse response.BookResponse
	copier.Copy(&bookResponse, book)
	user, err := controller.userService.GetUserById(int(book.CreatedBy))
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant get user",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	bookResponse.CreatedBy = user.FullName
	bookResponse.Shelve = book.Shelve.Name
	var tags []response.TagResponse
	var tagResponse response.TagResponse
	for _, t := range book.Tags {
		err := copier.Copy(&tagResponse, t)
		if err != nil {
			webResponse = response.WebResponse{
				Code:    http.StatusInternalServerError,
				Status:  "error",
				Message: "cant update book",
				Data:    nil,
			}
			c.JSON(http.StatusInternalServerError, webResponse)
			return
		}
		tags = append(tags, tagResponse)
	}
	bookResponse.Tags = tags
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "book updated",
		Data:    bookResponse,
	}
	c.JSON(http.StatusOK, webResponse)
}

// DeleteShelve godoc
// @Summary Delete a shelve
// @Description Delete a shelve by ID
// @Tags Shelve
// @Produce json
// @Param shelveId path int true "Shelve ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /shelves/{shelveId} [delete]
func (controller *BookController) DeleteShelve(c *gin.Context) {
	var webResponse response.WebResponse
	shelveIdStr := c.Param("shelveId")
	shelveId, err := strconv.Atoi(shelveIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get shelveId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.bookSerivce.DeleteShelve(shelveId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant delete shelve",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "shelve deleted",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}

// DeleteChapter godoc
// @Summary Delete a chapter
// @Description Delete a chapter by ID
// @Tags Chapter
// @Produce json
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /chapters/{chapterId} [delete]
func (controller *BookController) DeleteChapter(c *gin.Context) {
	var webResponse response.WebResponse
	chapterIdStr := c.Param("chapterId")
	chapterId, err := strconv.Atoi(chapterIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get chapterId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.bookSerivce.DeleteChapter(chapterId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant delete chapter",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "chapter deleted",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}

// DeletePage godoc
// @Summary Delete a page
// @Description Delete a page by ID
// @Tags Page
// @Produce json
// @Param pageId path int true "Page ID"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /chapter/{chapterId}/page/{pageId} [delete]
func (controller *BookController) DeletePage(c *gin.Context) {
	var webResponse response.WebResponse
	pageIdStr := c.Param("pageId")
	pageId, err := strconv.Atoi(pageIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "cant get pageId",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.bookSerivce.DeletePage(pageId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "cant delete page",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "page deleted",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}
