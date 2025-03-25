package repository

import (
	"bookstack/internal/models"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	CreatePermission(models.Permission) error
	GetPermissions() ([]models.Permission, error)
	DeletePermission(int) error
	FindIfExist(string) (*models.Permission, error)
	FindRoleBelong(string) ([]models.Role, error)
}

type PermissionRepositoryImpl struct {
	DB *gorm.DB
}

func NewPermissionRepositoryImpl(Db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{
		DB: Db,
	}
}
func (p *PermissionRepositoryImpl) FindIfExist(name string) (*models.Permission, error) {
	var permission models.Permission
	err := p.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (p *PermissionRepositoryImpl) FindRoleBelong(permission string) ([]models.Role, error) {
	var roles []models.Role
	result := p.DB.Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("permissions.name = ?", permission).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (p *PermissionRepositoryImpl) CreatePermission(permissionRequest models.Permission) error {
	err := p.DB.Create(&permissionRequest)
	if err != nil {
		return err.Error
	}
	return nil
}
func (p *PermissionRepositoryImpl) GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := p.DB.Find(&permissions)
	if err != nil {
		return nil, err.Error
	}
	return permissions, nil
}

func (p *PermissionRepositoryImpl) DeletePermission(permissionId int) error {
	// Giả sử bảng permissions có cột user_id lưu thông tin user
	result := p.DB.Delete(&models.Permission{}, permissionId)
	return result.Error
}
