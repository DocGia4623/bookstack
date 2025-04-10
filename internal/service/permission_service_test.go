package service

import (
	"bookstack/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

	tests := []struct {
		name        string
		permission  models.Permission
		mockError   error
		shouldError bool
		shouldCall  bool
	}{
		{
			name:        "successful creation with valid permission",
			permission:  models.Permission{Name: "test_permission"},
			mockError:   nil,
			shouldError: false,
			shouldCall:  true,
		},
		{
			name:        "empty permission name should still be created",
			permission:  models.Permission{Name: ""},
			mockError:   errors.New("permission name cannot be empty"),
			shouldError: true,
			shouldCall:  false,
		},
		{
			name:        "database error should be propagated",
			permission:  models.Permission{Name: "test_permission"},
			mockError:   errors.New("database error"),
			shouldError: true,
			shouldCall:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldCall {
				mockRepo.On("CreatePermission", tt.permission).Return(tt.mockError).Once()
			}

			err := service.CreatePermission(tt.permission)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.mockError.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldCall {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestGetPermissions(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionRepositoryImpl(mockRepo)

	tests := []struct {
		name            string
		mockPermissions []models.Permission
		mockError       error
		shouldError     bool
		expectedLength  int
	}{
		{
			name: "successfully retrieve multiple permissions",
			mockPermissions: []models.Permission{
				{Model: gorm.Model{ID: 1}, Name: "permission1"},
				{Model: gorm.Model{ID: 2}, Name: "permission2"},
				{Model: gorm.Model{ID: 3}, Name: "permission3"},
			},
			mockError:      nil,
			shouldError:    false,
			expectedLength: 3,
		},
		{
			name:            "successfully retrieve empty list",
			mockPermissions: []models.Permission{},
			mockError:       nil,
			shouldError:     false,
			expectedLength:  0,
		},
		{
			name:            "database error should return empty list",
			mockPermissions: []models.Permission{},
			mockError:       errors.New("database error"),
			shouldError:     true,
			expectedLength:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetPermissions").Return(tt.mockPermissions, tt.mockError).Once()

			permissions, err := service.GetPermissions()

			if tt.shouldError {
				assert.Error(t, err)
				assert.Empty(t, permissions)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, permissions)
				assert.Len(t, permissions, tt.expectedLength)

				// Verify each permission has required fields
				for _, p := range permissions {
					assert.NotEmpty(t, p.Name)
					assert.NotZero(t, p.ID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeletePermission(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionRepositoryImpl(mockRepo)

	tests := []struct {
		name         string
		permissionId int
		mockError    error
		shouldError  bool
		shouldCall   bool
	}{
		{
			name:         "successfully delete existing permission",
			permissionId: 1,
			mockError:    nil,
			shouldError:  false,
			shouldCall:   true,
		},
		{
			name:         "attempt to delete non-existent permission",
			permissionId: 999,
			mockError:    nil,
			shouldError:  false,
			shouldCall:   true,
		},
		{
			name:         "database error should be propagated",
			permissionId: 1,
			mockError:    errors.New("database error"),
			shouldError:  true,
			shouldCall:   true,
		},
		{
			name:         "invalid permission ID should not call repository",
			permissionId: 0,
			mockError:    nil,
			shouldError:  true,
			shouldCall:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldCall {
				mockRepo.On("DeletePermission", tt.permissionId).Return(tt.mockError).Once()
			}

			err := service.DeletePermission(tt.permissionId)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.mockError != nil {
					assert.Contains(t, err.Error(), tt.mockError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldCall {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
