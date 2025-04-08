package repository

import (
	"bookstack/internal/models"
	"regexp"
	"testing"

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
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPermissionRepositoryImpl(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "read", nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE name = $1 AND "permissions"."deleted_at" IS NULL ORDER BY "permissions"."id" LIMIT $2`)).
		WithArgs("read", 1).
		WillReturnRows(rows)

	permission, err := repo.FindIfExist("read")
	assert.NoError(t, err)
	assert.Equal(t, "read", permission.Name)
	assert.Equal(t, uint(1), permission.ID)
}

func TestCreatePermission(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPermissionRepositoryImpl(db)

	permission := models.Permission{
		Name: "write",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "permissions" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "write").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreatePermission(permission)
	assert.NoError(t, err)
}

func TestGetPermissions(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPermissionRepositoryImpl(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "read", nil, nil, nil).
		AddRow(2, "write", nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	permissions, err := repo.GetPermissions()
	assert.NoError(t, err)
	if assert.NotNil(t, permissions) {
		assert.Len(t, permissions, 2)
		assert.Equal(t, "read", permissions[0].Name)
		assert.Equal(t, "write", permissions[1].Name)
	}
}

func TestDeletePermission(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPermissionRepositoryImpl(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "permissions" SET "deleted_at"=$1 WHERE "permissions"."id" = $2 AND "permissions"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeletePermission(1)
	assert.NoError(t, err)
}
