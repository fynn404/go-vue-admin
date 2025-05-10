package handler

import (
	"net/http"
	"time"

	"github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/response"
	"github.com/fynn404/go-vue-admin/internal/service"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permService *service.PermissionService
}

func NewPermissionHandler(permService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permService: permService,
	}
}

// AssignRolePermissionRequest represents the request body for assigning a permission to a role
type AssignRolePermissionRequest struct {
	Role           string `json:"role" binding:"required"`
	PermissionCode string `json:"permission_code" binding:"required"`
}

// AssignUserPermissionRequest represents the request body for assigning a permission to a user
type AssignUserPermissionRequest struct {
	UserID         uint       `json:"user_id" binding:"required"`
	PermissionCode string     `json:"permission_code" binding:"required"`
	Reason         string     `json:"reason" binding:"required"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

// AssignRolePermission handles the request to assign a permission to a role
func (h *PermissionHandler) AssignRolePermission(c *gin.Context) {
	var req AssignRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Get the current user ID (the one performing the assignment)
	grantedBy, _ := c.Get("userID")

	err := h.permService.AssignRolePermission(
		c.Request.Context(),
		model.Role(req.Role),
		req.PermissionCode,
		grantedBy.(uint),
	)
	if err != nil {
		if err == service.ErrPermissionNotFound {
			response.NotFound(c, "Permission not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to assign permission", err)
		return
	}

	response.Success(c, nil)
}

// AssignUserPermission handles the request to assign a permission to a user
func (h *PermissionHandler) AssignUserPermission(c *gin.Context) {
	var req AssignUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Get the current user ID (the one performing the assignment)
	grantedBy, _ := c.Get("userID")

	err := h.permService.AssignUserPermission(
		c.Request.Context(),
		req.UserID,
		req.PermissionCode,
		grantedBy.(uint),
		req.Reason,
		req.ExpiresAt,
	)
	if err != nil {
		if err == service.ErrPermissionNotFound {
			response.NotFound(c, "Permission not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to assign permission", err)
		return
	}

	response.Success(c, nil)
}

// CheckPermission handles the request to check if a user has a specific permission
func (h *PermissionHandler) CheckPermission(c *gin.Context) {
	userID, _ := c.Get("userID")
	permissionCode := c.Query("permission_code")
	if permissionCode == "" {
		response.BadRequest(c, "Permission code is required")
		return
	}

	hasPermission, err := h.permService.HasPermission(
		c.Request.Context(),
		userID.(uint),
		permissionCode,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to check permission", err)
		return
	}

	response.Success(c, gin.H{
		"has_permission": hasPermission,
	})
}
