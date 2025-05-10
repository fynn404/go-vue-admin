package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/fynn404/go-vue-admin/internal/model"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// CreatePermission creates a new permission
func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetPermissionByCode gets a permission by its code
func (r *PermissionRepository) GetPermissionByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetUserPermissions gets all permissions for a user (including role permissions and user-specific permissions)
func (r *PermissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	// Get user's role permissions
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions_tab rp ON rp.permission_id = permissions_tab.id").
		Joins("JOIN users u ON u.role = rp.role").
		Where("u.id = ? AND permissions_tab.status = ? AND (rp.expires_at IS NULL OR rp.expires_at > ?)",
			userID, true, time.Now()).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	// Get user-specific permissions
	var userPermissions []*model.Permission
	err = r.db.WithContext(ctx).
		Joins("JOIN user_permissions_tab up ON up.permission_id = permissions_tab.id").
		Where("up.user_id = ? AND permissions_tab.status = ? AND (up.expires_at IS NULL OR up.expires_at > ?)",
			userID, true, time.Now()).
		Find(&userPermissions).Error
	if err != nil {
		return nil, err
	}

	// Merge permissions
	permissions = append(permissions, userPermissions...)
	return permissions, nil
}

// AssignRolePermission assigns a permission to a role
func (r *PermissionRepository) AssignRolePermission(ctx context.Context, rolePermission *model.RolePermission) error {
	return r.db.WithContext(ctx).Create(rolePermission).Error
}

// AssignUserPermission assigns a permission to a user
func (r *PermissionRepository) AssignUserPermission(ctx context.Context, userPermission *model.UserPermission) error {
	return r.db.WithContext(ctx).Create(userPermission).Error
}

// LogPermissionAudit logs a permission audit entry
func (r *PermissionRepository) LogPermissionAudit(ctx context.Context, audit *model.PermissionAudit) error {
	return r.db.WithContext(ctx).Create(audit).Error
}
