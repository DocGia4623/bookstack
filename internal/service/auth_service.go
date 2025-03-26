package service

import (
	"bookstack/config"
	"bookstack/helper"
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
	"bookstack/utils"
	"fmt"

	"github.com/jinzhu/copier"
)

type AuthService interface {
	Register(user request.UserCreateRequest) (models.User, error)
	Login(email, password string) (string, string, int, error)
	Logout(token string) error
	SaveRefreshToken(string, int) error
	RefreshToken(token string, signedKey string) (string, string, error)
}

type AuthServiceImpl struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewAuthServiceImpl(repo repository.UserRepository, conf *config.Config) AuthService {
	return &AuthServiceImpl{
		repo:   repo,
		config: conf,
	}
}
func (s *AuthServiceImpl) SaveRefreshToken(token string, userId int) error {
	return s.repo.SaveRefreshToken(token, userId)
}

func (s *AuthServiceImpl) Register(user request.UserCreateRequest) (models.User, error) {
	existingUser, _ := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return models.User{}, fmt.Errorf("user already exists")
	}
	var userModel models.User
	err := copier.Copy(&userModel, &user)
	if err != nil {
		return models.User{}, err
	}
	userModel.Password, err = utils.Hashpassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	userModel.EmailConfirmed = false
	userModel.RefreshToken = models.RefreshToken{
		Token: "",
	}
	responseUser, err := s.repo.NewUser(userModel)
	if err != nil {
		return models.User{}, err
	}
	return responseUser, nil
}
func (s *AuthServiceImpl) Login(email, password string) (string, string, int, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", 0, fmt.Errorf("user not found")
	}
	err = utils.VerifyPassword(user.Password, password)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid password")
	}
	accessToken, err := utils.GenerateAccessToken(s.config.AccessTokenExpiresIn, user.ID, s.config.AccessTokenSecret)
	if err != nil {
		return "", "", 0, err
	}
	// Generate refresh token
	refreshToken, err_refresh := utils.GenerateAccessToken(s.config.RefreshTokenExpiresIn, user.ID, s.config.RefreshTokenSecret)
	helper.ErrorPanic(err_refresh)
	return refreshToken, accessToken, user.ID, nil
}
func (s *AuthServiceImpl) Logout(token string) error {
	// err := utils.RevokeToken(token)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (a *AuthServiceImpl) RefreshToken(token string, signedKey string) (string, string, error) {
	// Kiểm tra refresh token có trong database không
	_, err := a.repo.FindByToken(token)
	if err != nil {
		return "", "", fmt.Errorf("cant find refresh token")
	}
	// Kiểm tra token có hợp lệ không
	config, _ := config.LoadConfig()
	sub, expRaw, err := utils.ValidateRefreshToken(token, config.RefreshTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("refresh token is invalid: %w", err)
	}

	// Xóa refresh token cũ khỏi database
	a.repo.DeleteToken(token)

	// Chuyển `expRaw` từ `interface{}` về `int64`
	expFloat, ok := expRaw.(float64)
	if !ok {
		return "", "", fmt.Errorf("invalid token expiration format")
	}
	exp := int64(expFloat) // Chuyển thành kiểu int64 (UNIX timestamp)

	// Tạo access token mới
	accessToken, err := utils.GenerateAccessToken(config.AccessTokenExpiresIn, sub, config.AccessTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate access token")
	}
	// Tạo refresh token mới
	newRefreshToken, err := utils.GenerateRefreshToken(exp, sub, config.RefreshTokenSecret)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate refresh token")
	}

	// Lưu refresh token mới vào database
	a.repo.SaveToken(models.RefreshToken{
		Token: newRefreshToken,
	})

	return accessToken, newRefreshToken, nil
}
