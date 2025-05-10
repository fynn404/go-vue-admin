package middleware

import (
	"net/http"

	"github.com/fynn404/go-vue-admin/internal/model"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Check if the user's role is in the allowed roles
		roleAllowed := false
		for _, role := range roles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware checks if the user has specific permissions
func PermissionMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// 定义角色权限映射
		rolePermissions := map[string][]string{
			string(model.RoleAdmin): {
				"manage_users", "manage_courses", "manage_system",
				"view_all_courses", "view_all_users", "manage_enrollments",
			},
			string(model.RoleTeacher): {
				"manage_own_courses", "view_own_courses", "manage_grades",
				"view_enrolled_students",
			},
			string(model.RoleStudent): {
				"view_available_courses", "enroll_courses", "view_own_grades",
				"view_own_schedule",
			},
		}

		// 检查用户是否有所需权限
		hasPermission := false
		if userPermissions, ok := rolePermissions[userRole.(string)]; ok {
			for _, required := range permissions {
				for _, granted := range userPermissions {
					if required == granted {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
