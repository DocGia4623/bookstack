package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Tags Users
// @Produce json
// @Success 200 {object} response.WebResponse "Successful retrieval of users"
// @Failure 500 {object} response.WebResponse "Service error"
// @Router /users [get]
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
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Get Users",
		Data:    users,
	}
	context.JSON(http.StatusOK, webResponse)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user details based on the token
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param user body request.UserUpdateRequest true "User update request"
// @Success 200 {object} response.WebResponse "Update successful"
// @Failure 400 {object} response.WebResponse "Invalid request or unauthorized"
// @Failure 500 {object} response.WebResponse "Failed to update user"
// @Router /users [put]
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

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Update user successful",
		Data:    user,
	}
	c.JSON(http.StatusOK, webResponse)
}
