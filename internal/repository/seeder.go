package repository

import (
	"bookstack/config"
	"bookstack/internal/constant"
	"bookstack/internal/models"
	"log"
)

// Seed roles & permissions
func SeedRolesAndPermissions() {
	db := config.DB

	// Define permissions
	permissions := []models.Permission{
		{Name: constant.ReadUser},
		{Name: constant.WriteUser},
		{Name: constant.DeleteUser},
		{Name: constant.ReceiveOrder},
		{Name: constant.UpdateOrderStatus},
	}

	// Tạo permissions
	for _, perm := range permissions {
		var existingPerm models.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			db.Create(&perm)
		}
	}

	// Define roles
	roles := []models.Role{
		{Name: "user"},
		{Name: "admin"},
		{Name: "editor"},
		{Name: "viewer"},
		{Name: "shipper"},
	}

	// Tạo roles
	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			db.Create(&role)
		}
	}

	// Lấy permissions cho admin
	var adminPermissions []models.Permission
	db.Where("name IN ?", []string{
		constant.ReadUser,
		constant.WriteUser,
		constant.DeleteUser,
	}).Find(&adminPermissions)

	// Lấy permissions cho shipper
	var shipperPermissions []models.Permission
	db.Where("name IN ?", []string{
		constant.ReceiveOrder,
	}).Find(&shipperPermissions)

	// Assign permissions to admin role
	var adminRole models.Role
	if err := db.First(&adminRole, "name = ?", "admin").Error; err == nil {
		// Xóa tất cả permissions cũ
		db.Model(&adminRole).Association("Permissions").Clear()
		// Thêm permissions mới
		db.Model(&adminRole).Association("Permissions").Append(adminPermissions)
	}

	// Assign permissions to shipper role
	var shipperRole models.Role
	if err := db.First(&shipperRole, "name = ?", "shipper").Error; err == nil {
		// Xóa tất cả permissions cũ
		db.Model(&shipperRole).Association("Permissions").Clear()
		// Thêm permissions mới
		db.Model(&shipperRole).Association("Permissions").Append(shipperPermissions)
	}

	log.Println("Seeded roles and permissions")
}
