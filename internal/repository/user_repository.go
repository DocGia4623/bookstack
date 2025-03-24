package repository

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserRepository interface {
	NewUser(models.User) (models.User, error)
	CreateUser(request.UserCreateRequest) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, user request.UserUpdateRequest) (models.User, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*models.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) NewUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) CreateUser(user request.UserCreateRequest) (models.User, error) {
	var userModel models.User
	err := copier.Copy(&userModel, &user)
	if err != nil {
		return models.User{}, err
	}
	return userModel, nil
}
func (r *UserRepositoryImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (r *UserRepositoryImpl) UpdateUser(id int, updateRequest request.UserUpdateRequest) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return models.User{}, err
	}

	err = copier.Copy(&user, &updateRequest)
	if err != nil {
		return models.User{}, err
	}

	err = r.db.Save(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
func (r *UserRepositoryImpl) DeleteUser(id int) error {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return err
	}

	err = r.db.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
