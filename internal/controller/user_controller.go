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

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// GetAllUser godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Authorization header string true "Authorization token"
// @Tags User
// @Produce json
// @Success 200 {object} response.WebResponse "Successful retrieval of users"
// @Failure 500 {object} response.WebResponse "Service error"
// @Router /user [get]
func (controller *UserController) GetAllUser(context *gin.Context) {
	var webResponse response.WebResponse
	users, err := controller.UserService.GetAllUsers()
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Fail",
			Message: "Service can't get user: " + err.Error(),
		}
		context.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var userResponse []response.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, controller.CoppyToUserResponse(user))
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Get Users",
		Data:    userResponse,
	}
	context.JSON(http.StatusOK, webResponse)
}

func (controller *UserController) CoppyToUserResponse(user models.User) response.UserResponse {
	var userResponse response.UserResponse
	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	return userResponse
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user details based on the token
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param user body request.UserUpdateRequest true "User update request"
// @Success 200 {object} response.WebResponse "Update successful"
// @Failure 400 {object} response.WebResponse "Invalid request or unauthorized"
// @Failure 500 {object} response.WebResponse "Failed to update user"
// @Router /user [put]
func (controller *UserController) UpdateUser(c *gin.Context) {
	var userUpdateRequest request.UserUpdateRequest
	var webResponse response.WebResponse
	header := c.Request.Header.Get("Authorization")

	userId, err := controller.UserService.GetUserIdByToken(header)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Fail",
			Message: "You are not logged in",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	if err := c.ShouldBindJSON(&userUpdateRequest); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Fail",
			Message: "Invalid request",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	user, err := controller.UserService.UpdateUser(userId, userUpdateRequest)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Fail",
			Message: "Failed to update: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	var userResponse response.UserResponse
	err = copier.Copy(&userResponse, user)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Fail",
			Message: "Failed to update: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Update user successful",
		Data:    userResponse,
	}
	c.JSON(http.StatusOK, webResponse)
}

// DeleteUser godoc
// @Summary Xóa người dùng
// @Description Xóa người dùng dựa trên ID
// @Tags User
// @Param userId path int true "ID của người dùng cần xóa"
// @Produce json
// @Success 200 {object} response.WebResponse "Xóa người dùng thành công"
// @Failure 400 {object} response.WebResponse "Không lấy được userId hợp lệ"
// @Failure 500 {object} response.WebResponse "Lỗi server"
// @Router /user/{userId} [delete]
func (controller *UserController) DeleteUser(c *gin.Context) {
	var webResponse response.WebResponse
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "Fail",
			Message: "Cant get userId",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.UserService.DeleteUser(userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "Fail",
			Message: "Server error" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Delete user successful",
	}
	c.JSON(http.StatusOK, webResponse)
}
