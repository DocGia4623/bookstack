package repository

import (
	"bookstack/config"
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"bookstack/utils"
	"errors"
	"fmt"

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
	GetUserById(int) (*models.User, error)
	FindIfUserHasRole(uint, []models.Role) error
	SaveRefreshToken(string, int) error
	SaveToken(refreshToken models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteToken(token string) error
	DeleteUserToken(userId int) error
}

type UserRepositoryImpl struct {
	db     *gorm.DB
	config *config.Config
}

func NewUserRepositoryImpl(db *gorm.DB, conf *config.Config) UserRepository {
	return &UserRepositoryImpl{
		db:     db,
		config: conf,
	}
}

func (u *UserRepositoryImpl) DeleteUserToken(userId int) error {
	result := u.db.Where("user_id = ?", userId).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserRepositoryImpl) SaveToken(refreshToken models.RefreshToken) error {
	result := u.db.Save(&refreshToken)
	return result.Error

}
func (u *UserRepositoryImpl) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	result := u.db.Where("token = ?", token).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}
func (u *UserRepositoryImpl) DeleteToken(token string) error {
	var refreshToken models.RefreshToken
	result := u.db.Where("token = ?", token).Delete(&refreshToken)
	return result.Error
}

func (u *UserRepositoryImpl) SaveRefreshToken(token string, userId int) error {
	user, err := u.GetUserById(userId)
	if err != nil {
		return err
	}
	if user.RefreshToken.Token != "" {
		_, _, err := utils.ValidateRefreshToken(user.RefreshToken.Token, u.config.RefreshTokenSecret)
		if err == nil {
			return nil
		}
	}
	// Xóa refresh token cũ của user (nếu có)
	if err := u.db.Where("user_id = ?", userId).Delete(&models.RefreshToken{}).Error; err != nil {
		return err
	}

	// Tạo refresh token mới
	newRefreshToken := models.RefreshToken{
		Token:  token,
		UserID: userId,
	}

	// Lưu refresh token mới vào database
	if err := u.db.Create(&newRefreshToken).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) GetUserById(userId int) (*models.User, error) {
	var existingUser models.User
	err := u.db.First(&existingUser, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &existingUser, nil
}

func (r *UserRepositoryImpl) GetRoleUser() (models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", "user").First(&role).Error
	return role, err
}

func (r *UserRepositoryImpl) NewUser(user models.User) (models.User, error) {
	role, err := r.GetRoleUser()
	if err != nil {
		return models.User{}, err
	}
	user.Roles = append(user.Roles, role)
	err = r.db.Create(&user).Error
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

	// Copy nhưng bỏ qua giá trị rỗng
	err = copier.CopyWithOption(&user, &updateRequest, copier.Option{IgnoreEmpty: true})
	if err != nil {
		return models.User{}, err
	}
	if updateRequest.Password != "" {
		password, err := utils.Hashpassword(updateRequest.Password)
		if err != nil {
			return models.User{}, err
		}
		user.Password = password
	}

	// Cập nhật dữ liệu
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

func (u *UserRepositoryImpl) FindIfUserHasRole(userID uint, roles []models.Role) error {
	var count int64

	// 🔹 Trích xuất tên role từ danh sách roles
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// 🔹 Truy vấn kiểm tra User có Role không
	result := u.db.Model(&models.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("users.id = ? AND roles.name IN ?", userID, roleNames).
		Select("COUNT(*)").Scan(&count)

	// 🔹 Kiểm tra lỗi query
	if result.Error != nil {
		return result.Error
	}

	// 🔹 Kiểm tra nếu không tìm thấy Role nào
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil // ✅ User có ít nhất một Role phù hợp
}
