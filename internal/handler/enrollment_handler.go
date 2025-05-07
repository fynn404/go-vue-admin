package handler

import (
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"

	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/model"
)

type UpdateGradeRequest struct {
	Grade decimal.Decimal `json:"grade" binding:"required"`
}

type BatchUpdateGradesRequest struct {
	Grades []struct {
		EnrollmentID uint            `json:"enrollment_id" binding:"required"`
		Grade        decimal.Decimal `json:"grade" binding:"required"`
	} `json:"grades" binding:"required"`
}

// EnrollmentHandler 处理选课相关的请求
type EnrollmentHandler struct{}

func NewEnrollmentHandler() *EnrollmentHandler {
	return &EnrollmentHandler{}
}

// List 获取选课列表
func (h *EnrollmentHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := config.DB.Model(&model.Enrollment{}).
		Preload("Course").
		Preload("Student")

	// Filter based on role
	switch role {
	case string(model.RoleStudent):
		query = query.Where("student_id = ?", userID)
	case string(model.RoleTeacher):
		query = query.Joins("JOIN courses ON enrollments.course_id = courses.id").
			Where("courses.teacher_id = ?", userID)
	}

	var enrollments []model.Enrollment
	var total int64

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := query.Offset(offset).Limit(limit).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": enrollments,
	})
}

// UpdateGrade 更新单个学生成绩
func (h *EnrollmentHandler) UpdateGrade(c *gin.Context) {
	enrollmentID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	var enrollment model.Enrollment
	if err := config.DB.Preload("Course").First(&enrollment, enrollmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	// 检查权限
	if enrollment.Course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var req UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := enrollment.UpdateGrade(req.Grade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, enrollment)
}

// GetCourseStats 获取课程统计信息
func (h *EnrollmentHandler) GetCourseStats(c *gin.Context) {
	courseID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	var course model.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// 检查权限
	if course.TeacherID != teacherID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var stats struct {
		TotalStudents int             `json:"total_students"`
		AvgGrade      decimal.Decimal `json:"avg_grade"`
		MaxGrade      decimal.Decimal `json:"max_grade"`
		MinGrade      decimal.Decimal `json:"min_grade"`
	}

	err := config.DB.Model(&model.Enrollment{}).
		Where("course_id = ? AND status = ?", courseID, model.EnrollmentStatusEnrolled).
		Select(`
			COUNT(*) as total_students,
			COALESCE(AVG(grade), 0) as avg_grade,
			COALESCE(MAX(grade), 0) as max_grade,
			COALESCE(MIN(grade), 0) as min_grade
		`).
		Scan(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetStudentGrades 获取学生成绩单
func (h *EnrollmentHandler) GetStudentGrades(c *gin.Context) {
	studentID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	// 只允许学生查看自己的成绩
	if role != string(model.RoleStudent) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var enrollments []model.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).
		Preload("Course").
		Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalCredits decimal.Decimal
	var totalGradePoints decimal.Decimal
	var completedCourses int

	for _, enrollment := range enrollments {
		if enrollment.Grade != nil && !enrollment.Grade.IsZero() {
			credits := decimal.NewFromInt(int64(enrollment.Course.Credits))
			totalCredits = totalCredits.Add(credits)

			// 计算绩点
			var points decimal.Decimal
			grade := enrollment.Grade.InexactFloat64()
			switch {
			case grade >= 90:
				points = decimal.NewFromFloat(4.0)
			case grade >= 85:
				points = decimal.NewFromFloat(3.7)
			case grade >= 80:
				points = decimal.NewFromFloat(3.3)
			case grade >= 75:
				points = decimal.NewFromFloat(3.0)
			case grade >= 70:
				points = decimal.NewFromFloat(2.7)
			case grade >= 65:
				points = decimal.NewFromFloat(2.3)
			case grade >= 60:
				points = decimal.NewFromFloat(2.0)
			default:
				points = decimal.Zero
			}

			totalGradePoints = totalGradePoints.Add(credits.Mul(points))
			completedCourses++
		}
	}

	var gpa decimal.Decimal
	if !totalCredits.IsZero() {
		gpa = totalGradePoints.Div(totalCredits)
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments":       enrollments,
		"total_credits":     totalCredits,
		"completed_courses": completedCourses,
		"gpa":               gpa.StringFixed(2),
	})
}

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

	var course model.Course
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
	var existingEnrollment model.Enrollment
	err := tx.Where("student_id = ? AND course_id = ? AND status = ?",
		studentID, courseID, model.EnrollmentStatusEnrolled).First(&existingEnrollment).Error
	if err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled in this course"})
		return
	}

	// 检查选课数量限制
	var activeEnrollments int64
	if err := tx.Model(&model.Enrollment{}).
		Where("student_id = ? AND status = ?", studentID, model.EnrollmentStatusEnrolled).
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
	enrollment := model.Enrollment{
		StudentID: studentID.(uint),
		CourseID:  course.ID,
		Status:    model.EnrollmentStatusEnrolled,
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

	var enrollment model.Enrollment
	err := config.DB.Where("student_id = ? AND course_id = ? AND status = ?",
		studentID, courseID, model.EnrollmentStatusEnrolled).First(&enrollment).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	// Use transaction to ensure data consistency
	tx := config.DB.Begin()

	// Update enrollment status
	enrollment.Status = model.EnrollmentStatusDropped
	if err := tx.Save(&enrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update enrollment"})
		return
	}

	// Update course enrollment count
	var course model.Course
	if err := tx.First(&course, courseID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find course"})
		return
	}

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

// GetCourseGradeStats returns grade statistics for a course
func GetCourseGradeStats(c *gin.Context) {
	courseID := c.Param("id")
	teacherID, _ := c.Get("user_id")

	var course model.Course
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
		TotalStudents int             `json:"total_students"`
		GradedCount   int             `json:"graded_count"`
		HighestGrade  decimal.Decimal `json:"highest_grade"`
		LowestGrade   decimal.Decimal `json:"lowest_grade"`
		AverageGrade  decimal.Decimal `json:"average_grade"`
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
		MaxGrade    decimal.Decimal
		MinGrade    decimal.Decimal
		AvgGrade    decimal.Decimal
	}

	config.DB.Model(&model.Enrollment{}).
		Where("course_id = ? AND status = ?", courseID, model.EnrollmentStatusEnrolled).
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
	var enrollments []model.Enrollment
	config.DB.Model(&model.Enrollment{}).
		Where("course_id = ? AND status = ? AND grade IS NOT NULL", courseID, model.EnrollmentStatusEnrolled).
		Find(&enrollments)

	ninety := decimal.NewFromInt(90)
	eighty := decimal.NewFromInt(80)
	seventy := decimal.NewFromInt(70)
	sixty := decimal.NewFromInt(60)

	for _, enrollment := range enrollments {
		if enrollment.Grade == nil {
			continue
		}
		switch {
		case enrollment.Grade.GreaterThanOrEqual(ninety):
			stats.Distribution.A++
		case enrollment.Grade.GreaterThanOrEqual(eighty):
			stats.Distribution.B++
		case enrollment.Grade.GreaterThanOrEqual(seventy):
			stats.Distribution.C++
		case enrollment.Grade.GreaterThanOrEqual(sixty):
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
	role, _ := c.Get("user_role")

	// 只允许学生查看自己的成绩
	if role != string(model.RoleStudent) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var enrollments []model.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).
		Preload("Course").
		Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalCredits decimal.Decimal
	var totalGradePoints decimal.Decimal
	var completedCourses int

	for _, enrollment := range enrollments {
		if enrollment.Grade != nil && !enrollment.Grade.IsZero() {
			credits := decimal.NewFromInt(int64(enrollment.Course.Credits))
			totalCredits = totalCredits.Add(credits)

			// 计算绩点
			var points decimal.Decimal
			grade := enrollment.Grade.InexactFloat64()
			switch {
			case grade >= 90:
				points = decimal.NewFromFloat(4.0)
			case grade >= 85:
				points = decimal.NewFromFloat(3.7)
			case grade >= 80:
				points = decimal.NewFromFloat(3.3)
			case grade >= 75:
				points = decimal.NewFromFloat(3.0)
			case grade >= 70:
				points = decimal.NewFromFloat(2.7)
			case grade >= 65:
				points = decimal.NewFromFloat(2.3)
			case grade >= 60:
				points = decimal.NewFromFloat(2.0)
			default:
				points = decimal.Zero
			}

			totalGradePoints = totalGradePoints.Add(credits.Mul(points))
			completedCourses++
		}
	}

	var gpa decimal.Decimal
	if !totalCredits.IsZero() {
		gpa = totalGradePoints.Div(totalCredits)
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments":       enrollments,
		"total_credits":     totalCredits,
		"completed_courses": completedCourses,
		"gpa":               gpa.StringFixed(2),
	})
}

// ListEnrollments returns a list of enrollments for a student or teacher
func ListEnrollments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := config.DB.Model(&model.Enrollment{}).
		Preload("Course").
		Preload("Student")

	// Filter based on role
	switch role {
	case string(model.RoleStudent):
		query = query.Where("student_id = ?", userID)
	case string(model.RoleTeacher):
		query = query.Joins("JOIN courses ON enrollments.course_id = courses.id").
			Where("courses.teacher_id = ?", userID)
	}

	var enrollments []model.Enrollment
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
//func GetRecommendedCourses(c *gin.Context) {
//	studentID, _ := c.Get("user_id")
//
//	// 从缓存获取推荐课程
//	key := fmt.Sprintf("recommendations:student:%d", studentID)
//	recommendedIDs, err := config.RedisClient.SMembers(config.Ctx, key).Result()
//	if err != nil {
//		// 如果缓存不存在，重新生成推荐
//		go updateCourseRecommendations(studentID.(uint))
//	}
//
//	var courses []model.Course
//	if len(recommendedIDs) > 0 {
//		if err := config.DB.Where("id IN ?", recommendedIDs).
//			Preload("Teacher").
//			Find(&courses).Error; err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended courses"})
//			return
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"courses": courses,
//	})
//}

// updateCourseRecommendations 更新学生的课程推荐
func updateCourseRecommendations(studentID uint) {
	// 获取学生已选课程
	var enrollments []model.Enrollment
	if err := config.DB.Where("student_id = ? AND status = ?",
		studentID, model.EnrollmentStatusEnrolled).
		Find(&enrollments).Error; err != nil {
		return
	}

	// 获取相似课程
	var courseIDs []uint
	for _, enrollment := range enrollments {
		courseIDs = append(courseIDs, enrollment.CourseID)
	}

	var recommendedCourses []model.Course
	if len(courseIDs) > 0 {
		// 基于已选课程推荐相似课程
		if err := config.DB.Where("id NOT IN ? AND status = ?",
			courseIDs, model.CourseStatusOpen).
			Where("credits IN (SELECT credits FROM courses WHERE id IN ?)", courseIDs).
			Or("teacher_id IN (SELECT teacher_id FROM courses WHERE id IN ?)", courseIDs).
			Limit(10).
			Find(&recommendedCourses).Error; err != nil {
			return
		}
	} else {
		// 如果没有选课历史，推荐热门课程
		if err := config.DB.Where("status = ?", model.CourseStatusOpen).
			Order("current_enrolled DESC").
			Limit(10).
			Find(&recommendedCourses).Error; err != nil {
			return
		}
	}

	// 更新推荐缓存
	//key := fmt.Sprintf("recommendations:student:%d", studentID)
	//pipe := config.RedisClient.Pipeline()
	//pipe.Del(config.Ctx, key)
	//for _, course := range recommendedCourses {
	//	pipe.SAdd(config.Ctx, key, course.ID)
	//}
	//pipe.Expire(config.Ctx, key, 24*time.Hour)
	//pipe.Exec(config.Ctx)
}
