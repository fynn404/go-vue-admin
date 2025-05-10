package service

import (
	"context"
	"errors"
	"time"

	"github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
)

var (
	ErrPermissionNotFound = errors.New("permission not found")
	ErrPermissionDenied   = errors.New("permission denied")
)

type PermissionService struct {
	permRepo *repository.PermissionRepository
}

func NewPermissionService(permRepo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		permRepo: permRepo,
	}
}

// HasPermission checks if a user has a specific permission
func (s *PermissionService) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	permissions, err := s.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p.Code == permissionCode {
			return true, nil
		}
	}
	return false, nil
}

// AssignRolePermission assigns a permission to a role
func (s *PermissionService) AssignRolePermission(ctx context.Context, role model.Role, permissionCode string, grantedBy uint) error {
	permission, err := s.permRepo.GetPermissionByCode(ctx, permissionCode)
	if err != nil {
		return err
	}

	rolePermission := &model.RolePermission{
		Role:         role,
		PermissionID: permission.ID,
		GrantedBy:    grantedBy,
		GrantedAt:    time.Now(),
	}

	if err := s.permRepo.AssignRolePermission(ctx, rolePermission); err != nil {
		return err
	}

	// Log the audit
	audit := &model.PermissionAudit{
		Role:          &role,
		PermissionID:  permission.ID,
		OperationType: "grant",
		OperatedBy:    grantedBy,
		OperatedAt:    time.Now(),
	}
	return s.permRepo.LogPermissionAudit(ctx, audit)
}

// AssignUserPermission assigns a permission to a user
func (s *PermissionService) AssignUserPermission(ctx context.Context, userID uint, permissionCode string, grantedBy uint, reason string, expiresAt *time.Time) error {
	permission, err := s.permRepo.GetPermissionByCode(ctx, permissionCode)
	if err != nil {
		return err
	}

	userPermission := &model.UserPermission{
		UserID:       userID,
		PermissionID: permission.ID,
		GrantedBy:    grantedBy,
		GrantedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		Reason:       reason,
	}

	if err := s.permRepo.AssignUserPermission(ctx, userPermission); err != nil {
		return err
	}

	// Log the audit
	audit := &model.PermissionAudit{
		UserID:        userID,
		PermissionID:  permission.ID,
		OperationType: "grant",
		OperatedBy:    grantedBy,
		OperatedAt:    time.Now(),
		Reason:        reason,
	}
	return s.permRepo.LogPermissionAudit(ctx, audit)
}
