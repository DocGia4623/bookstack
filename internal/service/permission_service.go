package service

import (
	"bookstack/internal/models"
	"bookstack/internal/repository"
)

type PermissionService interface {
	CreatePermission(models.Permission) error
	GetPermissions() ([]models.Permission, error)
	DeletePermission(int) error
}

type PermissionRepositoryImpl struct {
	repo repository.PermissionRepository
}

func NewPermissionRepositoryImpl(permissionRepo repository.PermissionRepository) PermissionService {
	return &PermissionRepositoryImpl{
		repo: permissionRepo,
	}
}
func (p *PermissionRepositoryImpl) CreatePermission(request models.Permission) error {
	return p.repo.CreatePermission(request)
}
func (p *PermissionRepositoryImpl) GetPermissions() ([]models.Permission, error) {
	return p.repo.GetPermissions()
}
func (p *PermissionRepositoryImpl) DeletePermission(permissionId int) error {
	return p.repo.DeletePermission(permissionId)
}
