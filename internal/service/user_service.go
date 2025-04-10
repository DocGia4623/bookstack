package service

import (
	"bookstack/config"
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/internal/repository"
	"bookstack/utils"
	"fmt"
	"strconv"
	"strings"
)

type UserService interface {
	CreateUser(request.UserCreateRequest) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, updateRequest request.UserUpdateRequest) (models.User, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserIdByToken(token string) (int, error)
	GetUserById(int) (*models.User, error)
	GetUserEmail(userId int) (string, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserServiceImpl(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) GetUserEmail(userId int) (string, error) {
	return s.repo.GetUserEmail(userId)
}

func (s *UserServiceImpl) GetUserById(userId int) (*models.User, error) {
	return s.repo.GetUserById(userId)
}

func (s *UserServiceImpl) GetUserIdByToken(token string) (int, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	config, err := config.LoadConfig()
	if err != nil {
		return 0, err
	}
	sub, err := utils.ValidateAccessToken(token, config.AccessTokenSecret)
	if err != nil {
		return 0, err
	}
	id, err_id := strconv.Atoi(fmt.Sprint(sub))
	if err_id != nil {
		return 0, err_id
	}
	return id, nil
}

func (s *UserServiceImpl) CreateUser(user request.UserCreateRequest) (models.User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}
func (s *UserServiceImpl) UpdateUser(id int, updateRequest request.UserUpdateRequest) (models.User, error) {
	return s.repo.UpdateUser(id, updateRequest)
}
func (s *UserServiceImpl) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
func (s *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}
