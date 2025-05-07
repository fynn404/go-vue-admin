package model

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;not null" json:"code"` // 权限代码，如 manage_users
	Name        string `gorm:"not null" json:"name"`             // 权限名称，如 "管理用户"
	Description string `json:"description"`                      // 权限描述
	Module      string `gorm:"not null" json:"module"`           // 所属模块，如 "user", "course"
}

// TableName 设置表名
func (Permission) TableName() string {
	return "permissions_tab"
}

// Role 角色类型
type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

// RolePermission 角色-权限关联模型
type RolePermission struct {
	gorm.Model
	Role         Role       `gorm:"not null" json:"role"`          // 角色
	PermissionID uint       `gorm:"not null" json:"permission_id"` // 权限ID
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "role_permissions_tab"
}
