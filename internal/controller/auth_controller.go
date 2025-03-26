package controller

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/dto/response"
	"bookstack/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	AuthenticationService service.AuthService
}

func NewAuthenticationController(authenticationService service.AuthService) *AuthenticationController {
	return &AuthenticationController{
		AuthenticationService: authenticationService,
	}
}

func (controller *AuthenticationController) Register(c *gin.Context) {
	var userRequest request.UserCreateRequest
	var webResponse response.WebResponse
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	user, err := controller.AuthenticationService.Register(userRequest)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}

	// // Gửi email xác nhận
	// err = utils.SendVerificationEmail(user.Email, user.FullName)
	// if err != nil {
	// 	webResponse = response.WebResponse{
	// 		Code:    http.StatusInternalServerError,
	// 		Status:  "error",
	// 		Message: "User registered but failed to send verification email",
	// 		Data:    nil,
	// 	}
	// 	c.JSON(http.StatusInternalServerError, webResponse)
	// 	return
	// }

	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User registered successfully, Please check your email to verify your account.",
		Data:    user,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *AuthenticationController) Login(c *gin.Context) {
	var userRequest request.UserLoginRequest
	var webResponse response.WebResponse
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	refreshToken, accessToken, userId, err := controller.AuthenticationService.Login(userRequest.Email, userRequest.Password)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)
	log.Println("Set refresh token in cookie: " + refreshToken) // ✅ Debug log

	err = controller.AuthenticationService.SaveRefreshToken(refreshToken, userId)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User logged in successfully",
		Data: response.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *AuthenticationController) Logout(c *gin.Context) {
	var webResponse response.WebResponse
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		webResponse = response.WebResponse{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Token not provided",
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, webResponse)
		return
	}
	err := controller.AuthenticationService.Logout(token)
	if err != nil {
		webResponse = response.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User logged out successfully",
		Data:    nil,
	}
	c.JSON(http.StatusOK, webResponse)
}
