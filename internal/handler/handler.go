// Package handler 处理HTTP请求的处理器
package handler

// Handler 包含所有HTTP处理器
type Handler struct {
	Auth                *AuthHandler
	User                *UserHandler
	Course              *CourseHandler
	Enrollment          *EnrollmentHandler
	Grade               *GradeHandler
	GradeHistoryHandler *GradeHistoryHandler
}

// NewHandler 创建一个新的Handler实例
func NewHandler() *Handler {
	return &Handler{
		Auth:                NewAuthHandler(),
		User:                NewUserHandler(),
		Course:              NewCourseHandler(),
		Enrollment:          NewEnrollmentHandler(),
		Grade:               NewGradeHandler(),
		GradeHistoryHandler: NewGradeHistoryHandler(),
	}
}
