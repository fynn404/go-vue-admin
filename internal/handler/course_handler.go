package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/model"
)

type CreateCourseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Credits     int    `json:"credits" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required"`
}

type UpdateCourseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int    `json:"credits"`
	Capacity    int    `json:"capacity"`
}

// CourseHandler 处理课程相关的请求
type CourseHandler struct{}

func NewCourseHandler() *CourseHandler {
	return &CourseHandler{}
}

// Create 创建新课程
func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID, _ := c.Get("user_id")
	course := model.Course{
		Name:        req.Name,
		Description: req.Description,
		TeacherID:   teacherID.(uint),
		Credits:     req.Credits,
		Capacity:    req.Capacity,
		Status:      model.CourseStatusOpen,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// Get 获取课程详情
func (h *CourseHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var course model.Course

	if err := config.DB.Preload("Teacher").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// List 获取课程列表
func (h *CourseHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	var courses []model.Course
	var total int64

	query := config.DB.Model(&model.Course{}).Preload("Teacher")

	// 应用搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 应用排序
	query = query.Order(sortBy + " " + sortOrder)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取分页数据
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 缓存热门课程
	if page == 1 && search == "" {
		go cacheHotCourses(courses)
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": courses,
	})
}

// cacheHotCourses 缓存热门课程信息
func cacheHotCourses(courses []model.Course) {
	for _, course := range courses {
		key := "course:" + strconv.FormatUint(uint64(course.ID), 10)
		data, _ := json.Marshal(course)
		config.Redis.Set(config.Ctx, key, data, 0)
	}
}

// Update 更新课程信息
func (h *CourseHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var course model.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 检查权限
	teacherID, _ := c.Get("user_id")
	if course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if req.Name != "" {
		course.Name = req.Name
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Credits != 0 {
		course.Credits = req.Credits
	}
	if req.Capacity != 0 {
		course.Capacity = req.Capacity
	}

	if err := config.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// Delete 删除课程
func (h *CourseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var course model.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 检查权限
	teacherID, _ := c.Get("user_id")
	if course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// Enroll 学生选课
func (h *CourseHandler) Enroll(c *gin.Context) {
	courseID := c.Param("id")
	studentID, _ := c.Get("user_id")

	var course model.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 检查课程状态
	if course.Status != model.CourseStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course is not open for enrollment"})
		return
	}

	// 检查是否已选
	var existingEnrollment model.Enrollment
	if err := config.DB.Where("course_id = ? AND student_id = ?", courseID, studentID).First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled"})
		return
	}

	// 检查容量
	var enrolledCount int64
	if err := config.DB.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&enrolledCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if int(enrolledCount) >= course.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course is full"})
		return
	}

	// 创建选课记录
	enrollment := model.Enrollment{
		CourseID:  course.ID,
		StudentID: studentID.(uint),
		Status:    model.EnrollmentStatusEnrolled,
	}

	if err := config.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

// Drop 退课
func (h *CourseHandler) Drop(c *gin.Context) {
	courseID := c.Param("id")
	studentID, _ := c.Get("user_id")

	var enrollment model.Enrollment
	if err := config.DB.Where("course_id = ? AND student_id = ?", courseID, studentID).First(&enrollment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	if err := config.DB.Delete(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course dropped successfully"})
}
