package middleware

import (
	"net/http"
	"strings"

	"go-vue-admin/config"
	"go-vue-admin/models"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the JWT token and sets the user in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user claims in the context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

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
			string(models.RoleAdmin): {
				"manage_users", "manage_courses", "manage_system",
				"view_all_courses", "view_all_users", "manage_enrollments",
			},
			string(models.RoleTeacher): {
				"manage_own_courses", "view_own_courses", "manage_grades",
				"view_enrolled_students",
			},
			string(models.RoleStudent): {
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

// ResourceOwnerMiddleware checks if the user owns the resource
func ResourceOwnerMiddleware(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("role")
		resourceID := c.Param("id")

		// 管理员可以访问所有资源
		if userRole == string(models.RoleAdmin) {
			c.Next()
			return
		}

		var isOwner bool
		switch resourceType {
		case "course":
			var course models.Course
			if err := config.DB.First(&course, resourceID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
				c.Abort()
				return
			}
			isOwner = course.TeacherID == userID.(uint)

		case "enrollment":
			var enrollment models.Enrollment
			if err := config.DB.First(&enrollment, resourceID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
				c.Abort()
				return
			}
			isOwner = enrollment.StudentID == userID.(uint)
		}

		if !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
