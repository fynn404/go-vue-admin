package handler

import (
	"net/http"
	"strconv"

	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/model"

	"github.com/fynn404/go-vue-admin/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// GradeRequest 成绩请求结构
type GradeRequest struct {
	Score   decimal.Decimal `json:"score" binding:"required"`
	Comment string          `json:"comment"`
}

// GradeHandler 处理成绩相关的请求
type GradeHandler struct {
	gradeService service.GradeService
}

func NewGradeHandler() *GradeHandler {
	return &GradeHandler{
		gradeService: service.NewGradeService(),
	}
}

// handleError handles common error responses
func (h *GradeHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrCourseNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
	case service.ErrStudentNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
	case service.ErrNotAuthorized:
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
	case service.ErrGradeAlreadyExist:
		c.JSON(http.StatusConflict, gin.H{"error": "Grade already exists"})
	case service.ErrGradeNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
	case service.ErrInvalidScore:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
	case service.ErrInvalidCourseID:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
	case service.ErrInvalidStudentID:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Create 创建成绩
func (h *GradeHandler) Create(c *gin.Context) {
	var req GradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID, _ := c.Get("user_id")
	courseID := c.Param("course_id")
	studentID := c.Param("student_id")

	grade, err := h.gradeService.CreateGradeWithValidation(courseID, studentID, teacherID.(uint), req.Score, req.Comment)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, grade)
}

// Update 更新成绩
func (h *GradeHandler) Update(c *gin.Context) {
	var req GradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gradeID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	grade, err := h.gradeService.UpdateGradeWithValidation(gradeID, teacherID.(uint), req.Score, req.Comment)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, grade)
}

// Publish 发布成绩
func (h *GradeHandler) Publish(c *gin.Context) {
	gradeID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	gradeIDUint, err := strconv.ParseUint(gradeID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	grade, err := h.gradeService.PublishGrade(uint(gradeIDUint), teacherID.(uint))
	if err != nil {
		switch err {
		case service.ErrGradeNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
		case service.ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish grade"})
		}
		return
	}

	c.JSON(http.StatusOK, grade)
}

// GetStudentGrades 获取学生成绩列表
func (h *GradeHandler) GetStudentGrades(c *gin.Context) {
	studentID, _ := c.Get("user_id")

	grades, err := h.gradeService.GetStudentGrades(studentID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grades"})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetCourseGrades 获取课程成绩列表
func (h *GradeHandler) GetCourseGrades(c *gin.Context) {
	teacherID, _ := c.Get("user_id")
	courseID := c.Param("course_id")

	courseIDUint, err := strconv.ParseUint(courseID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	grades, err := h.gradeService.GetCourseGrades(uint(courseIDUint), teacherID.(uint))
	if err != nil {
		switch err {
		case service.ErrCourseNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		case service.ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grades"})
		}
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetGradeHistory 获取成绩修改历史
func (h *GradeHandler) GetGradeHistory(c *gin.Context) {
	gradeID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	gradeIDUint, err := strconv.ParseUint(gradeID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	histories, err := h.gradeService.GetGradeHistory(uint(gradeIDUint), userID.(uint), role.(string))
	if err != nil {
		switch err {
		case service.ErrGradeNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
		case service.ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grade history"})
		}
		return
	}

	c.JSON(http.StatusOK, histories)
}

// BatchUpdateGrades 批量更新成绩
func (h *GradeHandler) BatchUpdateGrades(c *gin.Context) {
	teacherID, _ := c.Get("user_id")

	var req BatchUpdateGradesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, gradeUpdate := range req.Grades {
		var enrollment model.Enrollment
		if err := tx.Preload("Course").First(&enrollment, gradeUpdate.EnrollmentID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
			return
		}

		// 检查权限
		if enrollment.Course.TeacherID != teacherID.(uint) {
			tx.Rollback()
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}

		if err := enrollment.UpdateGrade(gradeUpdate.Grade); err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := tx.Save(&enrollment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grades updated successfully"})
}
