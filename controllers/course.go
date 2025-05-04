package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-vue-admin/config"
	"go-vue-admin/models"
)

type CreateCourseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Credits     float32 `json:"credits" binding:"required"`
	Capacity    int     `json:"capacity" binding:"required"`
}

// CreateCourse creates a new course
func CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID, _ := c.Get("user_id")
	course := models.Course{
		Name:        req.Name,
		Description: req.Description,
		TeacherID:   teacherID.(uint),
		Credits:     req.Credits,
		Capacity:    req.Capacity,
		Status:      models.CourseStatusOpen,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}

// GetCourse returns a specific course
func GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.Preload("Teacher").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// ListCourses returns a list of courses with pagination and filters
func ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 获取筛选参数
	search := c.Query("search")
	status := c.Query("status")
	teacherID := c.Query("teacher_id")
	credits := c.Query("credits")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	var courses []models.Course
	var total int64

	query := config.DB.Model(&models.Course{}).Preload("Teacher")

	// 应用搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 应用筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if credits != "" {
		query = query.Where("credits = ?", credits)
	}

	// 应用排序
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count courses"})
		return
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	// 使用 Redis 缓存热门课程
	if page == 1 && search == "" {
		go cacheHotCourses(courses)
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// cacheHotCourses 缓存热门课程信息
func cacheHotCourses(courses []models.Course) {
	for _, course := range courses {
		key := fmt.Sprintf("course:%d", course.ID)
		data, err := json.Marshal(course)
		if err != nil {
			continue
		}
		config.RedisClient.Set(config.Ctx, key, data, 30*time.Minute)
	}
}

// SearchCourses handles course search with advanced filters
func SearchCourses(c *gin.Context) {
	var req struct {
		Search    string   `json:"search"`
		Status    []string `json:"status"`
		Credits   []string `json:"credits"`
		TeacherID uint     `json:"teacher_id"`
		Page      int      `json:"page"`
		Limit     int      `json:"limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit
	query := config.DB.Model(&models.Course{}).Preload("Teacher")

	// 构建搜索条件
	if req.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	}

	if len(req.Credits) > 0 {
		query = query.Where("credits IN ?", req.Credits)
	}

	if req.TeacherID > 0 {
		query = query.Where("teacher_id = ?", req.TeacherID)
	}

	var total int64
	var courses []models.Course

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count courses"})
		return
	}

	if err := query.Offset(offset).Limit(req.Limit).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
		"total":   total,
		"page":    req.Page,
		"limit":   req.Limit,
	})
}

// UpdateCourse updates a course
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Check if the user is the teacher of this course
	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own courses"})
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.Name = req.Name
	course.Description = req.Description
	course.Credits = req.Credits
	course.Capacity = req.Capacity

	if err := config.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
		"course":  course,
	})
}

// DeleteCourse deletes a course
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Check if the user is the teacher of this course
	userID, _ := c.Get("user_id")
	if course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own courses"})
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course deleted successfully",
	})
}
