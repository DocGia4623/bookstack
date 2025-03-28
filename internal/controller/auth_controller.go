package controller

import (
	"bookstack/config"
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

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserCreateRequest true "User registration data"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 500 {object} response.WebResponse
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Logs in a user and returns access & refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserLoginRequest true "Login credentials"
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Failure 401 {object} response.WebResponse
// @Router /auth/login [post]
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
			TokenType:    "Bearer token",
		},
	}
	c.JSON(http.StatusOK, webResponse)
}

// Logout godoc
// @Summary User logout
// @Description Logs out a user by invalidating their token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.WebResponse
// @Failure 401 {object} response.WebResponse
// @Router /auth/logout [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access token using a refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.WebResponse
// @Failure 400 {object} response.WebResponse
// @Router /auth/refresh [post]
func (controller *AuthenticationController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.WebResponse{Code: http.StatusUnauthorized, Status: "unauthorized", Message: "Refresh token is missing"})
		return
	}
	log.Println("refresh Token: " + refreshToken)
	config, _ := config.LoadConfig()
	newAccessToken, newRefreshToken, err := controller.AuthenticationService.RefreshToken(refreshToken, config.RefreshTokenSecret)
	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.WebResponse{Code: http.StatusBadRequest, Status: "bad request", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.WebResponse{Code: http.StatusOK, Status: "ok", Message: "Refresh token success", Data: response.LoginResponse{TokenType: "Bearer Token", RefreshToken: newRefreshToken, AccessToken: newAccessToken}})
}
