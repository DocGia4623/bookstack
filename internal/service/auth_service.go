package service

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
	"bookstack/utils"
	"fmt"

	"github.com/jinzhu/copier"
)

type AuthService interface {
	Register(user request.UserCreateRequest) (models.User, error)
	Login(email, password string) (string, error)
	Logout(token string) error
}

type AuthServiceImpl struct {
	repo repository.UserRepository
}

func NewAuthServiceImpl(repo repository.UserRepository) AuthService {
	return &AuthServiceImpl{
		repo: repo,
	}
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
	responseUser, err := s.repo.NewUser(userModel)
	if err != nil {
		return models.User{}, err
	}
	return responseUser, nil
}
func (s *AuthServiceImpl) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	err = utils.VerifyPassword(user.Password, password)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	// token, err := utils.GenerateToken(user.ID)
	// if err != nil {
	// 	return "", err
	// }
	return "Token not done", nil
}
func (s *AuthServiceImpl) Logout(token string) error {
	// err := utils.RevokeToken(token)
	// if err != nil {
	// 	return err
	// }
	return nil
}
