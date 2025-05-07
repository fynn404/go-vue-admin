package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fynn404/go-vue-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// GradeHistoryHandler 处理成绩历史记录相关的请求
type GradeHistoryHandler struct {
	historyService service.GradeHistoryService
}

// NewGradeHistoryHandler 创建成绩历史记录处理器实例
func NewGradeHistoryHandler() *GradeHistoryHandler {
	return &GradeHistoryHandler{
		historyService: service.NewGradeHistoryService(),
	}
}

// GetHistoryByGrade 获取指定成绩的历史记录
func (h *GradeHistoryHandler) GetHistoryByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")

	// 转换成绩ID
	gradeIDUint, err := strconv.ParseUint(gradeID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	histories, err := h.historyService.GetHistoryByGradeID(uint(gradeIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grade history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// GetHistoryByStudent 获取当前学生的成绩历史记录
func (h *GradeHistoryHandler) GetHistoryByStudent(c *gin.Context) {
	userID, _ := c.Get("user_id")

	histories, err := h.historyService.GetHistoryByStudentID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// GetHistoryByCourse 获取指定课程的成绩历史记录
func (h *GradeHistoryHandler) GetHistoryByCourse(c *gin.Context) {
	courseID := c.Param("course_id")

	// 转换课程ID
	courseIDUint, err := strconv.ParseUint(courseID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	histories, err := h.historyService.GetHistoryByCourseID(uint(courseIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// GetHistoryByTeacher 获取当前教师的操作历史记录
func (h *GradeHistoryHandler) GetHistoryByTeacher(c *gin.Context) {
	userID, _ := c.Get("user_id")

	histories, err := h.historyService.GetHistoryByTeacherID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teacher history"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// GetHistoryDetail 获取历史记录详情
func (h *GradeHistoryHandler) GetHistoryDetail(c *gin.Context) {
	historyID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// 转换历史记录ID
	historyIDUint, err := strconv.ParseUint(historyID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid history ID"})
		return
	}

	history, err := h.historyService.GetHistoryDetail(uint(historyIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History record not found"})
		return
	}

	// 权限检查：只有管理员、相关教师和学生可以查看
	if role != "admin" && userID.(uint) != history.StudentID && userID.(uint) != history.TeacherID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetHistoryByDateRange 获取指定日期范围内的历史记录
func (h *GradeHistoryHandler) GetHistoryByDateRange(c *gin.Context) {
	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// 设置结束日期为当天的最后一刻
	endDate = endDate.Add(24*time.Hour - time.Second)

	histories, err := h.historyService.GetHistoryByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history records"})
		return
	}

	c.JSON(http.StatusOK, histories)
}
