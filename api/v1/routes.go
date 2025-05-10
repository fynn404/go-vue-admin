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
	auth         *AuthRoutes
	user         *UserRoutes
	course       *CourseRoutes
	enrollment   *EnrollmentRoutes
	grade        *GradeRoutes
	gradeHistory *GradeHistoryRoutes
}

// NewRoutes 创建路由实例
func NewRoutes(h *controllers.Handler) *Routes {
	return &Routes{
		auth:         NewAuthRoutes(h),
		user:         NewUserRoutes(h),
		course:       NewCourseRoutes(h),
		enrollment:   NewEnrollmentRoutes(h),
		grade:        NewGradeRoutes(h),
		gradeHistory: NewGradeHistoryRoutes(h),
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

		// 教师路由 - 使用细粒度权限控制
		teacher := courses.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			createCourse := teacher.Group("", middleware.PermissionMiddleware("manage_own_courses"))
			{
				createCourse.POST("", r.handler.Course.Create)
				createCourse.PUT("/:id", r.handler.Course.Update)
				createCourse.DELETE("/:id", r.handler.Course.Delete)
			}
		}

		// 学生路由 - 使用细粒度权限控制
		student := courses.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			enrollCourse := student.Group("", middleware.PermissionMiddleware("enroll_courses"))
			{
				enrollCourse.POST("/:id/enroll", r.handler.Course.Enroll)
				enrollCourse.POST("/:id/drop", r.handler.Course.Drop)
			}
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

		// 教师路由 - 使用细粒度权限控制
		teacher := enrollments.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			manageGrades := teacher.Group("", middleware.PermissionMiddleware("manage_grades"))
			{
				manageGrades.PUT("/:id/grade", r.handler.Enrollment.UpdateGrade)
				manageGrades.GET("/courses/:id/stats", r.handler.Enrollment.GetCourseStats)
			}
		}

		// 学生路由 - 使用细粒度权限控制
		student := enrollments.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			viewGrades := student.Group("", middleware.PermissionMiddleware("view_own_grades"))
			{
				viewGrades.GET("/grades", r.handler.Enrollment.GetStudentGrades)
			}
		}
	}
}

// GradeRoutes 成绩相关路由
type GradeRoutes struct {
	handler *controllers.Handler
}

func NewGradeRoutes(h *controllers.Handler) *GradeRoutes {
	return &GradeRoutes{handler: h}
}

func (r *GradeRoutes) Register(group *gin.RouterGroup) {
	grades := group.Group("/grades", middleware.AuthMiddleware())
	{
		// 学生路由 - 使用细粒度权限控制
		student := grades.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			viewGrades := student.Group("", middleware.PermissionMiddleware("view_own_grades"))
			{
				viewGrades.GET("/my", r.handler.Grade.GetStudentGrades)
			}
		}

		// 教师路由 - 使用细粒度权限控制
		teacher := grades.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			manageGrades := teacher.Group("", middleware.PermissionMiddleware("manage_grades"))
			{
				manageGrades.POST("/courses/:course_id/students/:student_id", r.handler.Grade.Create)
				manageGrades.GET("/courses/:course_id", r.handler.Grade.GetCourseGrades)
				manageGrades.PUT("/:id", r.handler.Grade.Update)
				manageGrades.POST("/:id/publish", r.handler.Grade.Publish)
			}
		}

		// 通用路由（需要权限验证）
		grades.GET("/:id/history", r.handler.Grade.GetGradeHistory)
	}
}

// GradeHistoryRoutes 成绩历史记录相关路由
type GradeHistoryRoutes struct {
	handler *controllers.Handler
}

func NewGradeHistoryRoutes(h *controllers.Handler) *GradeHistoryRoutes {
	return &GradeHistoryRoutes{handler: h}
}

func (r *GradeHistoryRoutes) Register(group *gin.RouterGroup) {
	histories := group.Group("/grade-histories", middleware.AuthMiddleware())
	{
		// 学生路由 - 使用细粒度权限控制
		student := histories.Group("", middleware.RoleMiddleware(string(model.RoleStudent)))
		{
			viewGrades := student.Group("", middleware.PermissionMiddleware("view_own_grades"))
			{
				viewGrades.GET("/my", r.handler.GradeHistoryHandler.GetHistoryByStudent)
			}
		}

		// 教师路由 - 使用细粒度权限控制
		teacher := histories.Group("", middleware.RoleMiddleware(string(model.RoleTeacher)))
		{
			manageGrades := teacher.Group("", middleware.PermissionMiddleware("manage_grades"))
			{
				manageGrades.GET("/courses/:course_id", r.handler.GradeHistoryHandler.GetHistoryByCourse)
				manageGrades.GET("/grades/:grade_id", r.handler.GradeHistoryHandler.GetHistoryByGrade)
				manageGrades.GET("/my-operations", r.handler.GradeHistoryHandler.GetHistoryByTeacher)
			}
		}

		// 管理员路由 - 使用细粒度权限控制
		admin := histories.Group("", middleware.RoleMiddleware(string(model.RoleAdmin)))
		{
			manageSystem := admin.Group("", middleware.PermissionMiddleware("manage_system"))
			{
				manageSystem.GET("/date-range", r.handler.GradeHistoryHandler.GetHistoryByDateRange)
			}
		}

		// 通用路由（需要权限验证）
		histories.GET("/:id", r.handler.GradeHistoryHandler.GetHistoryDetail)
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
	routes.grade.Register(v1)
	routes.gradeHistory.Register(v1)
}
