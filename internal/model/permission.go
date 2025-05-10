package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;not null" json:"code"` // 权限代码，如 manage_users
	Name        string `gorm:"not null" json:"name"`             // 权限名称，如 "管理用户"
	Description string `json:"description"`                      // 权限描述
	Module      string `gorm:"not null" json:"module"`           // 所属模块，如 "user", "course"
	Type        string `gorm:"not null" json:"type"`             // 权限类型：menu-菜单权限，operation-操作权限
	Status      bool   `gorm:"default:true" json:"status"`       // 权限状态：true-启用，false-禁用
}

// TableName 设置表名
func (Permission) TableName() string {
	return "permissions_tab"
}

// RolePermission 角色-权限关联模型
type RolePermission struct {
	gorm.Model
	Role         Role       `gorm:"not null" json:"role"`          // 角色
	PermissionID uint       `gorm:"not null" json:"permission_id"` // 权限ID
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
	GrantedBy    uint       `gorm:"not null" json:"granted_by"` // 授权人ID
	GrantedAt    time.Time  `gorm:"not null" json:"granted_at"` // 授权时间
	ExpiresAt    *time.Time `json:"expires_at"`                 // 过期时间，为空表示永不过期
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "role_permissions_tab"
}

// UserPermission 用户-权限关联模型（用于特殊权限分配）
type UserPermission struct {
	gorm.Model
	UserID       uint       `gorm:"not null" json:"user_id"`       // 用户ID
	PermissionID uint       `gorm:"not null" json:"permission_id"` // 权限ID
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
	GrantedBy    uint       `gorm:"not null" json:"granted_by"` // 授权人ID
	GrantedAt    time.Time  `gorm:"not null" json:"granted_at"` // 授权时间
	ExpiresAt    *time.Time `json:"expires_at"`                 // 过期时间，为空表示永不过期
	Reason       string     `json:"reason"`                     // 授权原因
}

// TableName 设置表名
func (UserPermission) TableName() string {
	return "user_permissions_tab"
}

// PermissionAudit 权限变更审计日志
type PermissionAudit struct {
	gorm.Model
	UserID        uint      `gorm:"not null" json:"user_id"`        // 被授权的用户ID
	Role          *Role     `json:"role"`                           // 被授权的角色，为空表示是用户特殊权限
	PermissionID  uint      `gorm:"not null" json:"permission_id"`  // 权限ID
	OperationType string    `gorm:"not null" json:"operation_type"` // 操作类型：grant-授权，revoke-撤销
	OperatedBy    uint      `gorm:"not null" json:"operated_by"`    // 操作人ID
	OperatedAt    time.Time `gorm:"not null" json:"operated_at"`    // 操作时间
	Reason        string    `json:"reason"`                         // 操作原因
}

// TableName 设置表名
func (PermissionAudit) TableName() string {
	return "permission_audits_tab"
}
