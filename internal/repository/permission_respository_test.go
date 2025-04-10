package repository

import (
	"bookstack/internal/models"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return db, mock, cleanup
}

func TestFindIfExist(t *testing.T) {
	tests := []struct {
		name           string
		permissionName string
		mockRows       *sqlmock.Rows
		expectedError  error
		expectFound    bool
	}{
		{
			name:           "permission exists",
			permissionName: "read",
			mockRows: sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "read", time.Now(), time.Now(), nil),
			expectedError: nil,
			expectFound:   true,
		},
		{
			name:           "permission does not exist",
			permissionName: "nonexistent",
			mockRows:       sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}),
			expectedError:  nil,
			expectFound:    false,
		},
		{
			name:           "database error",
			permissionName: "read",
			mockRows:       nil,
			expectedError:  errors.New("database error"),
			expectFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewPermissionRepositoryImpl(db)

			query := regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE name = $1 AND "permissions"."deleted_at" IS NULL ORDER BY "permissions"."id" LIMIT $2`)

			if tt.expectedError != nil {
				mock.ExpectQuery(query).
					WithArgs(tt.permissionName, 1).
					WillReturnError(tt.expectedError)
			} else {
				mock.ExpectQuery(query).
					WithArgs(tt.permissionName, 1).
					WillReturnRows(tt.mockRows)
			}

			permission, err := repo.FindIfExist(tt.permissionName)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectFound {
					assert.NotNil(t, permission)
					assert.Equal(t, tt.permissionName, permission.Name)
				} else {
					assert.Nil(t, permission)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatePermission(t *testing.T) {
	tests := []struct {
		name          string
		permission    models.Permission
		expectedError error
	}{
		{
			name: "successful creation",
			permission: models.Permission{
				Name: "write",
			},
			expectedError: nil,
		},
		{
			name: "database error",
			permission: models.Permission{
				Name: "error_case",
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewPermissionRepositoryImpl(db)

			if tt.expectedError == nil {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "permissions" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, tt.permission.Name).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			} else {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "permissions"`)).
					WillReturnError(tt.expectedError)
				mock.ExpectRollback()
			}

			err := repo.CreatePermission(tt.permission)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPermissions(t *testing.T) {
	tests := []struct {
		name           string
		mockRows       *sqlmock.Rows
		expectedError  error
		expectedLength int
		expectedNames  []string
	}{
		{
			name: "successful retrieval multiple permissions",
			mockRows: sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "read", time.Now(), time.Now(), nil).
				AddRow(2, "write", time.Now(), time.Now(), nil).
				AddRow(3, "delete", time.Now(), time.Now(), nil),
			expectedError:  nil,
			expectedLength: 3,
			expectedNames:  []string{"read", "write", "delete"},
		},
		{
			name:           "empty result",
			mockRows:       sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}),
			expectedError:  nil,
			expectedLength: 0,
			expectedNames:  []string{},
		},
		{
			name:           "database error",
			mockRows:       nil,
			expectedError:  errors.New("database error"),
			expectedLength: 0,
			expectedNames:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewPermissionRepositoryImpl(db)

			query := regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."deleted_at" IS NULL`)

			if tt.expectedError != nil {
				mock.ExpectQuery(query).WillReturnError(tt.expectedError)
			} else {
				mock.ExpectQuery(query).WillReturnRows(tt.mockRows)
			}

			permissions, err := repo.GetPermissions()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, permissions)
			} else {
				assert.NoError(t, err)
				assert.Len(t, permissions, tt.expectedLength)
				if tt.expectedLength > 0 {
					for i, name := range tt.expectedNames {
						assert.Equal(t, name, permissions[i].Name)
					}
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePermission(t *testing.T) {
	tests := []struct {
		name          string
		permissionID  int
		expectedError error
		rowsAffected  int64
	}{
		{
			name:          "successful deletion",
			permissionID:  1,
			expectedError: nil,
			rowsAffected:  1,
		},
		{
			name:          "permission not found",
			permissionID:  999,
			expectedError: nil,
			rowsAffected:  0,
		},
		{
			name:          "database error",
			permissionID:  1,
			expectedError: errors.New("database error"),
			rowsAffected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewPermissionRepositoryImpl(db)

			if tt.expectedError != nil {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "permissions" SET "deleted_at"=$1`)).
					WithArgs(sqlmock.AnyArg(), tt.permissionID).
					WillReturnError(tt.expectedError)
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "permissions" SET "deleted_at"=$1 WHERE "permissions"."id" = $2 AND "permissions"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), tt.permissionID).
					WillReturnResult(sqlmock.NewResult(1, tt.rowsAffected))
				mock.ExpectCommit()
			}

			err := repo.DeletePermission(tt.permissionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
