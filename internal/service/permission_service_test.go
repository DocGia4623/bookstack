package service

import (
	"bookstack/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) CreatePermission(permission models.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetPermissions() ([]models.Permission, error) {
	args := m.Called()
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) DeletePermission(permissionId int) error {
	args := m.Called(permissionId)
	return args.Error(0)
}

func (m *MockPermissionRepository) FindIfExist(name string) (*models.Permission, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindRoleBelong(permissionId string) ([]models.Role, error) {
	args := m.Called(permissionId)
	return args.Get(0).([]models.Role), args.Error(1)
}

func TestCreatePermission(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionRepositoryImpl(mockRepo)

	permission := models.Permission{
		Name: "test_permission",
	}

	// Test successful creation
	mockRepo.On("CreatePermission", permission).Return(nil).Once()
	err := service.CreatePermission(permission)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test error case
	expectedError := errors.New("database error")
	mockRepo.On("CreatePermission", permission).Return(expectedError).Once()
	err = service.CreatePermission(permission)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetPermissions(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionRepositoryImpl(mockRepo)

	// Test successful retrieval
	expectedPermissions := []models.Permission{
		{Name: "permission1"},
		{Name: "permission2"},
	}
	mockRepo.On("GetPermissions").Return(expectedPermissions, nil).Once()
	permissions, err := service.GetPermissions()
	assert.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)
	mockRepo.AssertExpectations(t)

	// Test error case
	expectedError := errors.New("database error")
	mockRepo.On("GetPermissions").Return([]models.Permission{}, expectedError).Once()
	permissions, err = service.GetPermissions()
	assert.Equal(t, expectedError, err)
	assert.Empty(t, permissions)
	mockRepo.AssertExpectations(t)
}

func TestDeletePermission(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionRepositoryImpl(mockRepo)

	permissionId := 1

	// Test successful deletion
	mockRepo.On("DeletePermission", permissionId).Return(nil).Once()
	err := service.DeletePermission(permissionId)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test error case
	expectedError := errors.New("database error")
	mockRepo.On("DeletePermission", permissionId).Return(expectedError).Once()
	err = service.DeletePermission(permissionId)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
