package v1

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/fynn404/go-vue-admin/internal/handler"
	"github.com/fynn404/go-vue-admin/internal/middleware"
	"github.com/fynn404/go-vue-admin/internal/model"
)

// RouteGroup 定义路由组接口
type RouteGroup interface {
	Register(group *gin.RouterGroup)
}

// Routes 路由分组结构体
type Routes struct {
	auth       *AuthRoutes
	user       *UserRoutes
	course     *CourseRoutes
	enrollment *EnrollmentRoutes
}

// NewRoutes 创建路由实例
func NewRoutes(h *controllers.Handler) *Routes {
	return &Routes{
		auth:       NewAuthRoutes(h),
		user:       NewUserRoutes(h),
		course:     NewCourseRoutes(h),
		enrollment: NewEnrollmentRoutes(h),
	}
}

// AuthRoutes 认证相关路由
type AuthRoutes struct {
	handler *controllers.Handler
}

func NewAuthRoutes(h *controllers.Handler) *AuthRoutes {
	return &AuthRoutes{handler: h}
}

func (r *AuthRoutes) Register(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	{
		auth.POST("/login", r.handler.Auth.Login)
		auth.POST("/register", r.handler.Auth.Register)
	}
}

// UserRoutes 用户相关路由
type UserRoutes struct {
	handler *controllers.Handler
}

func NewUserRoutes(h *controllers.Handler) *UserRoutes {
	return &UserRoutes{handler: h}
}

func (r *UserRoutes) Register(group *gin.RouterGroup) {
	users := group.Group("/users", middleware.AuthMiddleware())
	{
		users.GET("/profile", r.handler.User.GetProfile)
		users.PUT("/profile", r.handler.User.UpdateProfile)
	}
}

// CourseRoutes 课程相关路由
type CourseRoutes struct {
	handler *controllers.Handler
}

func NewCourseRoutes(h *controllers.Handler) *CourseRoutes {
	return &CourseRoutes{handler: h}
}

func (r *CourseRoutes) Register(group *gin.RouterGroup) {
	courses := group.Group("/courses", middleware.AuthMiddleware())
	{
		// 公共路由
		courses.GET("", r.handler.Course.List)
		courses.GET("/:id", r.handler.Course.Get)

		// 教师路由
		teacher := courses.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			teacher.POST("", r.handler.Course.Create)
			teacher.PUT("/:id", r.handler.Course.Update)
			teacher.DELETE("/:id", r.handler.Course.Delete)
		}

		// 学生路由
		student := courses.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			student.POST("/:id/enroll", r.handler.Course.Enroll)
			student.POST("/:id/drop", r.handler.Course.Drop)
		}
	}
}

// EnrollmentRoutes 选课相关路由
type EnrollmentRoutes struct {
	handler *controllers.Handler
}

func NewEnrollmentRoutes(h *controllers.Handler) *EnrollmentRoutes {
	return &EnrollmentRoutes{handler: h}
}

func (r *EnrollmentRoutes) Register(group *gin.RouterGroup) {
	enrollments := group.Group("/enrollments", middleware.AuthMiddleware())
	{
		enrollments.GET("", r.handler.Enrollment.List)

		// 教师路由
		teacher := enrollments.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			teacher.PUT("/:id/grade", r.handler.Enrollment.UpdateGrade)
			teacher.POST("/grades/batch", r.handler.Enrollment.BatchUpdateGrades)
			teacher.GET("/courses/:id/stats", r.handler.Enrollment.GetCourseStats)
		}

		// 学生路由
		student := enrollments.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			student.GET("/grades", r.handler.Enrollment.GetStudentGrades)
		}
	}
}

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine, h *controllers.Handler) {
	// API v1 版本分组
	v1 := r.Group("/api/v1")

	// 创建路由实例
	routes := NewRoutes(h)

	// 注册各个模块的路由
	routes.auth.Register(v1)
	routes.user.Register(v1)
	routes.course.Register(v1)
	routes.enrollment.Register(v1)
}
