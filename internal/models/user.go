package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             int    `json:"id"`
	FullName       string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RememberToken  string `json:"remember_token"`
	EmailConfirmed bool   `json:"email_confirmed"`
	ImageId        int    `json:"image_id"`
	Roles          []Role `gorm:"many2many:user_roles"`
}

// Role struct
type Role struct {
	gorm.Model
	ID          int          `json:"id"`
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

// Permission struct
type Permission struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `gorm:"unique"`
}
