package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-vue-admin/config"
	"go-vue-admin/models"
)

// EnrollCourse handles course enrollment for students
func EnrollCourse(c *gin.Context) {
	courseID := c.Param("id")
	studentID, _ := c.Get("user_id")

	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var course models.Course
	if err := tx.First(&course, courseID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 检查课程是否可选
	if !course.IsAvailable() {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course is not available for enrollment"})
		return
	}

	// 检查是否已经选过这门课
	var existingEnrollment models.Enrollment
	err := tx.Where("student_id = ? AND course_id = ? AND status = ?",
		studentID, courseID, models.EnrollmentStatusActive).First(&existingEnrollment).Error
	if err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled in this course"})
		return
	}

	// 检查选课数量限制
	var activeEnrollments int64
	if err := tx.Model(&models.Enrollment{}).
		Where("student_id = ? AND status = ?", studentID, models.EnrollmentStatusActive).
		Count(&activeEnrollments).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment limit"})
		return
	}

	if activeEnrollments >= 6 { // 假设每个学期最多选6门课
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum enrollment limit reached"})
		return
	}

	// 创建选课记录
	enrollment := models.Enrollment{
		StudentID: studentID.(uint),
		CourseID:  course.ID,
		Status:    models.EnrollmentStatusActive,
	}

	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enrollment"})
		return
	}

	// 更新课程已选人数
	course.CurrentEnrolled++
	if err := tx.Save(&course).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// 异步更新课程推荐缓存
	go updateCourseRecommendations(studentID.(uint))

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Successfully enrolled in course",
		"enrollment": enrollment,
	})
}

// DropCourse handles course dropping for students
func DropCourse(c *gin.Context) {
	courseID := c.Param("id")
	studentID, _ := c.Get("user_id")

	var enrollment models.Enrollment
	err := config.DB.Where("student_id = ? AND course_id = ? AND status = ?",
		studentID, courseID, models.EnrollmentStatusActive).First(&enrollment).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	// Use transaction to ensure data consistency
	tx := config.DB.Begin()

	// Update enrollment status
	enrollment.Status = models.EnrollmentStatusDropped
	if err := tx.Save(&enrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update enrollment"})
		return
	}

	// Update course enrollment count
	var course models.Course
	if err := tx.First(&course, courseID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find course"})
		return
	}

	course.CurrentEnrolled--
	if err := tx.Save(&course).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully dropped course",
	})
}

// UpdateGrade handles grade updates for teachers
func UpdateGrade(c *gin.Context) {
	enrollmentID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	var enrollment models.Enrollment
	if err := config.DB.Preload("Course").First(&enrollment, enrollmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	// 验证教师权限
	if enrollment.Course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only grade your own courses"})
		return
	}

	var req struct {
		Grade     float32 `json:"grade" binding:"required,min=0,max=100"`
		Comment   string  `json:"comment"`
		Feedback  string  `json:"feedback"`
		Timestamp string  `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开启事务
	tx := config.DB.Begin()

	// 更新成绩
	enrollment.Grade = &req.Grade
	if err := tx.Save(&enrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grade"})
		return
	}

	// 记录成绩变更历史
	gradeHistory := models.GradeHistory{
		EnrollmentID: enrollment.ID,
		Grade:        req.Grade,
		Comment:      req.Comment,
		TeacherID:    teacherID.(uint),
		Timestamp:    time.Now(),
	}

	if err := tx.Create(&gradeHistory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record grade history"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Grade updated successfully",
		"enrollment": enrollment,
	})
}

// BatchUpdateGrades handles batch grade updates for teachers
func BatchUpdateGrades(c *gin.Context) {
	teacherID, _ := c.Get("user_id")

	var req struct {
		Grades []struct {
			EnrollmentID uint    `json:"enrollment_id" binding:"required"`
			Grade        float32 `json:"grade" binding:"required,min=0,max=100"`
			Comment      string  `json:"comment"`
		} `json:"grades" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开启事务
	tx := config.DB.Begin()

	for _, gradeUpdate := range req.Grades {
		var enrollment models.Enrollment
		if err := tx.Preload("Course").First(&enrollment, gradeUpdate.EnrollmentID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Enrollment %d not found", gradeUpdate.EnrollmentID)})
			return
		}

		// 验证教师权限
		if enrollment.Course.TeacherID != teacherID.(uint) {
			tx.Rollback()
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only grade your own courses"})
			return
		}

		// 更新成绩
		enrollment.Grade = &gradeUpdate.Grade
		if err := tx.Save(&enrollment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grades"})
			return
		}

		// 记录成绩变更历史
		gradeHistory := models.GradeHistory{
			EnrollmentID: enrollment.ID,
			Grade:        gradeUpdate.Grade,
			Comment:      gradeUpdate.Comment,
			TeacherID:    teacherID.(uint),
			Timestamp:    time.Now(),
		}

		if err := tx.Create(&gradeHistory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record grade history"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Grades updated successfully",
	})
}

// GetCourseGradeStats returns grade statistics for a course
func GetCourseGradeStats(c *gin.Context) {
	courseID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 验证教师权限
	if course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view statistics for your own courses"})
		return
	}

	// 获取成绩统计
	var stats struct {
		TotalStudents int     `json:"total_students"`
		GradedCount   int     `json:"graded_count"`
		HighestGrade  float32 `json:"highest_grade"`
		LowestGrade   float32 `json:"lowest_grade"`
		AverageGrade  float32 `json:"average_grade"`
		Distribution  struct {
			A int `json:"a"` // 90-100
			B int `json:"b"` // 80-89
			C int `json:"c"` // 70-79
			D int `json:"d"` // 60-69
			F int `json:"f"` // 0-59
		} `json:"distribution"`
	}

	// 计算总体统计
	var result struct {
		Count       int
		GradedCount int
		MaxGrade    float32
		MinGrade    float32
		AvgGrade    float32
	}

	config.DB.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ?", courseID, models.EnrollmentStatusActive).
		Select(`
			COUNT(*) as count,
			COUNT(grade) as graded_count,
			MAX(grade) as max_grade,
			MIN(grade) as min_grade,
			AVG(grade) as avg_grade
		`).
		Scan(&result)

	stats.TotalStudents = result.Count
	stats.GradedCount = result.GradedCount
	stats.HighestGrade = result.MaxGrade
	stats.LowestGrade = result.MinGrade
	stats.AverageGrade = result.AvgGrade

	// 计算成绩分布
	var grades []float32
	config.DB.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ? AND grade IS NOT NULL", courseID, models.EnrollmentStatusActive).
		Pluck("grade", &grades)

	for _, grade := range grades {
		switch {
		case grade >= 90:
			stats.Distribution.A++
		case grade >= 80:
			stats.Distribution.B++
		case grade >= 70:
			stats.Distribution.C++
		case grade >= 60:
			stats.Distribution.D++
		default:
			stats.Distribution.F++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetStudentGrades returns all grades for a student
func GetStudentGrades(c *gin.Context) {
	studentID, _ := c.Get("user_id")

	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ? AND grade IS NOT NULL", studentID).
		Preload("Course").
		Preload("Course.Teacher").
		Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
		return
	}

	// 计算GPA
	var totalPoints float32
	var totalCredits float32
	for _, enrollment := range enrollments {
		if enrollment.Grade != nil {
			grade := *enrollment.Grade
			credits := enrollment.Course.Credits

			// GPA计算规则
			var points float32
			switch {
			case grade >= 90:
				points = 4.0
			case grade >= 85:
				points = 3.7
			case grade >= 80:
				points = 3.3
			case grade >= 75:
				points = 3.0
			case grade >= 70:
				points = 2.7
			case grade >= 65:
				points = 2.3
			case grade >= 60:
				points = 2.0
			default:
				points = 0.0
			}

			totalPoints += points * credits
			totalCredits += credits
		}
	}

	var gpa float32
	if totalCredits > 0 {
		gpa = totalPoints / totalCredits
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments":   enrollments,
		"gpa":           gpa,
		"total_credits": totalCredits,
	})
}

// ListEnrollments returns a list of enrollments for a student or teacher
func ListEnrollments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Enrollment{}).
		Preload("Course").
		Preload("Student")

	// Filter based on role
	switch role {
	case string(models.RoleStudent):
		query = query.Where("student_id = ?", userID)
	case string(models.RoleTeacher):
		query = query.Joins("JOIN courses ON enrollments.course_id = courses.id").
			Where("courses.teacher_id = ?", userID)
	}

	var enrollments []models.Enrollment
	var total int64

	// Get total count
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// GetRecommendedCourses returns recommended courses for a student
func GetRecommendedCourses(c *gin.Context) {
	studentID, _ := c.Get("user_id")

	// 从缓存获取推荐课程
	key := fmt.Sprintf("recommendations:student:%d", studentID)
	recommendedIDs, err := config.RedisClient.SMembers(config.Ctx, key).Result()
	if err != nil {
		// 如果缓存不存在，重新生成推荐
		go updateCourseRecommendations(studentID.(uint))
	}

	var courses []models.Course
	if len(recommendedIDs) > 0 {
		if err := config.DB.Where("id IN ?", recommendedIDs).
			Preload("Teacher").
			Find(&courses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended courses"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

// updateCourseRecommendations 更新学生的课程推荐
func updateCourseRecommendations(studentID uint) {
	// 获取学生已选课程
	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ? AND status = ?",
		studentID, models.EnrollmentStatusActive).
		Find(&enrollments).Error; err != nil {
		return
	}

	// 获取相似课程
	var courseIDs []uint
	for _, enrollment := range enrollments {
		courseIDs = append(courseIDs, enrollment.CourseID)
	}

	var recommendedCourses []models.Course
	if len(courseIDs) > 0 {
		// 基于已选课程推荐相似课程
		if err := config.DB.Where("id NOT IN ? AND status = ?",
			courseIDs, models.CourseStatusOpen).
			Where("credits IN (SELECT credits FROM courses WHERE id IN ?)", courseIDs).
			Or("teacher_id IN (SELECT teacher_id FROM courses WHERE id IN ?)", courseIDs).
			Limit(10).
			Find(&recommendedCourses).Error; err != nil {
			return
		}
	} else {
		// 如果没有选课历史，推荐热门课程
		if err := config.DB.Where("status = ?", models.CourseStatusOpen).
			Order("current_enrolled DESC").
			Limit(10).
			Find(&recommendedCourses).Error; err != nil {
			return
		}
	}

	// 更新推荐缓存
	key := fmt.Sprintf("recommendations:student:%d", studentID)
	pipe := config.RedisClient.Pipeline()
	pipe.Del(config.Ctx, key)
	for _, course := range recommendedCourses {
		pipe.SAdd(config.Ctx, key, course.ID)
	}
	pipe.Expire(config.Ctx, key, 24*time.Hour)
	pipe.Exec(config.Ctx)
}
