package routes

import (
	"github.com/gin-gonic/gin"

	"go-vue-admin/controllers"
	"go-vue-admin/middleware"
	"go-vue-admin/models"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine) {
	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			user := protected.Group("/users")
			{
				user.GET("/profile", controllers.GetProfile)
				user.PUT("/profile", controllers.UpdateProfile)
			}

			// Course routes
			courses := protected.Group("/courses")
			{
				courses.GET("", controllers.ListCourses)
				courses.GET("/:id", controllers.GetCourse)

				// Teacher only routes
				teacherOnly := courses.Group("")
				teacherOnly.Use(middleware.RoleMiddleware(string(models.RoleTeacher)))
				{
					teacherOnly.POST("", controllers.CreateCourse)
					teacherOnly.PUT("/:id", controllers.UpdateCourse)
					teacherOnly.DELETE("/:id", controllers.DeleteCourse)
				}

				// Student only routes
				studentOnly := courses.Group("")
				studentOnly.Use(middleware.RoleMiddleware(string(models.RoleStudent)))
				{
					studentOnly.POST("/:id/enroll", controllers.EnrollCourse)
					studentOnly.POST("/:id/drop", controllers.DropCourse)
				}
			}

			// Enrollment routes
			enrollments := protected.Group("/enrollments")
			{
				enrollments.GET("", controllers.ListEnrollments)

				// Teacher only routes
				teacherOnly := enrollments.Group("")
				teacherOnly.Use(middleware.RoleMiddleware(string(models.RoleTeacher)))
				{
					teacherOnly.PUT("/:id/grade", controllers.UpdateGrade)
					teacherOnly.POST("/grades/batch", controllers.BatchUpdateGrades)
					teacherOnly.GET("/courses/:id/stats", controllers.GetCourseGradeStats)
				}

				// Student only routes
				studentOnly := enrollments.Group("")
				studentOnly.Use(middleware.RoleMiddleware(string(models.RoleStudent)))
				{
					studentOnly.GET("/grades", controllers.GetStudentGrades)
				}
			}
		}
	}
}
