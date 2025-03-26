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
