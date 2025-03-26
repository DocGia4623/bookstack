package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             int          `json:"id"`
	FullName       string       `json:"name"`
	Email          string       `json:"email"`
	Password       string       `json:"password"`
	RememberToken  string       `json:"remember_token"`
	EmailConfirmed bool         `json:"email_confirmed"`
	ImageId        int          `json:"image_id"`
	Roles          []Role       `gorm:"many2many:user_roles"`
	RefreshToken   RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Role struct
type Role struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type RolePermission struct {
	RoleID       uint
	PermissionID uint
}

// Permission struct
type Permission struct {
	gorm.Model
	Name string `gorm:"unique"`
}

//refreshToken
type RefreshToken struct {
	ID     uint   `json:"ID"`
	Token  string `json:"token"`
	UserID int    `gorm:"unique;not null"`
}
