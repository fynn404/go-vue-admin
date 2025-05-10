package middleware

import (
	"context"
	"net/http"

	"github.com/fynn404/go-vue-admin/internal/response"
	"github.com/fynn404/go-vue-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// RequirePermission creates a middleware that checks if the user has the required permission
func RequirePermission(permSvc *service.PermissionService, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by authentication middleware)
		userID, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		// Check permission
		hasPermission, err := permSvc.HasPermission(c.Request.Context(), userID.(uint), permissionCode)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to check permission", err)
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by authentication middleware)
		userRole, exists := c.Get("userRole")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		if userRole.(string) != role {
			response.Forbidden(c, "Role not authorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

// ResourceOwner creates a middleware that checks if the user owns the resource
func ResourceOwner(resourceType string, resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		resourceID := c.Param(resourceIDParam)
		if resourceID == "" {
			response.BadRequest(c, "Resource ID not provided")
			c.Abort()
			return
		}

		// Here you would implement the logic to check if the user owns the resource
		// This is just a placeholder - you need to implement the actual check based on your data model
		isOwner := checkResourceOwnership(c.Request.Context(), resourceType, resourceID, userID.(uint))
		if !isOwner {
			response.Forbidden(c, "Not resource owner")
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkResourceOwnership is a placeholder function - implement based on your data model
func checkResourceOwnership(ctx context.Context, resourceType, resourceID string, userID uint) bool {
	// Implement the actual ownership check here
	// This might involve querying your database to check if the user owns the resource
	return false
}
