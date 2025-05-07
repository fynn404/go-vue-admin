package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fynn404/go-vue-admin/internal/config"
	models "github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrEnrollmentNotFound   = errors.New("enrollment not found")
	ErrCourseNotAvailable   = errors.New("course is not available for enrollment")
	ErrAlreadyEnrolled      = errors.New("already enrolled in this course")
	ErrMaxEnrollmentReached = errors.New("maximum enrollment limit reached")
	ErrUnauthorized         = errors.New("unauthorized access")
)

// EnrollmentService 定义选课服务接口
type EnrollmentService interface {
	// EnrollCourse 学生选课
	EnrollCourse(ctx context.Context, studentID, courseID uint) error
	// DropCourse 退课
	DropCourse(ctx context.Context, enrollmentID uint) error
	// UpdateGrade 更新成绩
	UpdateGrade(ctx context.Context, enrollmentID uint, grade decimal.Decimal) error
	// GetEnrollment 获取选课记录
	GetEnrollment(ctx context.Context, enrollmentID uint) (*models.Enrollment, error)
	// ListStudentEnrollments 获取学生的选课列表
	ListStudentEnrollments(ctx context.Context, studentID uint, page, pageSize int) ([]*models.Enrollment, int64, error)
	// ListCourseEnrollments 获取课程的选课列表
	ListCourseEnrollments(ctx context.Context, courseID uint, page, pageSize int) ([]*models.Enrollment, int64, error)
	// List 获取选课列表
	List(userID uint, role string, offset, limit int) ([]models.Enrollment, int64, error)
	// GetStudentGrades 获取学生成绩单
	GetStudentGrades(studentID uint) (*struct {
		Enrollments      []models.Enrollment `json:"enrollments"`
		TotalCredits     decimal.Decimal     `json:"total_credits"`
		CompletedCourses int                 `json:"completed_courses"`
		GPA              string              `json:"gpa"`
	}, error)
	// GetCourseGradeStats 获取课程成绩统计
	GetCourseGradeStats(courseID, teacherID uint) (*struct {
		TotalStudents int             `json:"total_students"`
		GradedCount   int             `json:"graded_count"`
		HighestGrade  decimal.Decimal `json:"highest_grade"`
		LowestGrade   decimal.Decimal `json:"lowest_grade"`
		AverageGrade  decimal.Decimal `json:"average_grade"`
		Distribution  struct {
			A int `json:"a"`
			B int `json:"b"`
			C int `json:"c"`
			D int `json:"d"`
			F int `json:"f"`
		} `json:"distribution"`
	}, error)
}

// EnrollmentServiceImpl 实现 EnrollmentService 接口
type EnrollmentServiceImpl struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
	db             *gorm.DB
}

// NewEnrollmentService 创建选课服务实例
func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	userRepo repository.UserRepository,
) EnrollmentService {
	return &EnrollmentServiceImpl{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
		db:             config.DB,
	}
}

// EnrollCourse 实现学生选课
func (s *EnrollmentServiceImpl) EnrollCourse(ctx context.Context, studentID, courseID uint) error {
	// 验证学生身份
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return models.ErrInvalidStudent
	}
	if student.Role != models.RoleStudent {
		return models.ErrInvalidStudent
	}

	// 验证课程是否存在并可选
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return models.ErrCourseUnavailable
	}
	if !course.IsAvailable() {
		return models.ErrCourseUnavailable
	}

	// 检查是否已经选过这门课
	var existingEnrollment models.Enrollment
	err = s.db.Where("student_id = ? AND course_id = ? AND status = ?",
		studentID, courseID, models.EnrollmentStatusEnrolled).First(&existingEnrollment).Error
	if err == nil {
		return ErrAlreadyEnrolled
	}

	// 检查选课数量限制
	var activeEnrollments int64
	err = s.db.Model(&models.Enrollment{}).
		Where("student_id = ? AND status = ?", studentID, models.EnrollmentStatusEnrolled).
		Count(&activeEnrollments).Error
	if err != nil {
		return err
	}

	if activeEnrollments >= 6 { // 假设每个学期最多选6门课
		return ErrMaxEnrollmentReached
	}

	// 创建选课记录
	enrollment := &models.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    models.EnrollmentStatusEnrolled,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return err
	}

	// 更新课程选课人数
	course.CurrentEnrolled++
	return s.courseRepo.Update(ctx, course)
}

// DropCourse 实现退课
func (s *EnrollmentServiceImpl) DropCourse(ctx context.Context, enrollmentID uint) error {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return ErrEnrollmentNotFound
	}

	if err := enrollment.Drop(); err != nil {
		return err
	}

	// 更新选课记录
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return err
	}

	// 更新课程选课人数
	course, err := s.courseRepo.FindByID(ctx, enrollment.CourseID)
	if err != nil {
		return err
	}

	// 获取学生信息
	student, err := s.userRepo.FindByID(ctx, enrollment.StudentID)
	if err != nil {
		return err
	}

	if err := course.Unenroll(*student); err != nil {
		return err
	}

	return s.courseRepo.Update(ctx, course)
}

// UpdateGrade 实现更新成绩
func (s *EnrollmentServiceImpl) UpdateGrade(ctx context.Context, enrollmentID uint, grade decimal.Decimal) error {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return ErrEnrollmentNotFound
	}

	if err := enrollment.UpdateGrade(grade); err != nil {
		return models.ErrInvalidGrade
	}

	return s.enrollmentRepo.Update(ctx, enrollment)
}

// GetEnrollment 实现获取选课记录
func (s *EnrollmentServiceImpl) GetEnrollment(ctx context.Context, enrollmentID uint) (*models.Enrollment, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}
	return enrollment, nil
}

// ListStudentEnrollments 实现获取学生的选课列表
func (s *EnrollmentServiceImpl) ListStudentEnrollments(ctx context.Context, studentID uint, page, pageSize int) ([]*models.Enrollment, int64, error) {
	offset := (page - 1) * pageSize
	return s.enrollmentRepo.ListByStudent(ctx, studentID, offset, pageSize)
}

// ListCourseEnrollments 实现获取课程的选课列表
func (s *EnrollmentServiceImpl) ListCourseEnrollments(ctx context.Context, courseID uint, page, pageSize int) ([]*models.Enrollment, int64, error) {
	offset := (page - 1) * pageSize
	return s.enrollmentRepo.ListByCourse(ctx, courseID, offset, pageSize)
}

// List 实现获取选课列表
func (s *EnrollmentServiceImpl) List(userID uint, role string, offset, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64
	var err error

	switch role {
	case "admin":
		// 管理员可以查看所有选课记录
		err = s.db.Model(&models.Enrollment{}).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = s.db.Offset(offset).Limit(limit).Find(&enrollments).Error
	case "teacher":
		// 教师只能查看自己课程的选课记录
		err = s.db.Model(&models.Enrollment{}).
			Joins("JOIN courses ON enrollments.course_id = courses.id").
			Where("courses.teacher_id = ?", userID).
			Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = s.db.Joins("JOIN courses ON enrollments.course_id = courses.id").
			Where("courses.teacher_id = ?", userID).
			Offset(offset).Limit(limit).
			Find(&enrollments).Error
	case "student":
		// 学生只能查看自己的选课记录
		err = s.db.Model(&models.Enrollment{}).
			Where("student_id = ?", userID).
			Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = s.db.Where("student_id = ?", userID).
			Offset(offset).Limit(limit).
			Find(&enrollments).Error
	default:
		return nil, 0, ErrUnauthorized
	}

	if err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

// GetStudentGrades 实现获取学生成绩单
func (s *EnrollmentServiceImpl) GetStudentGrades(studentID uint) (*struct {
	Enrollments      []models.Enrollment `json:"enrollments"`
	TotalCredits     decimal.Decimal     `json:"total_credits"`
	CompletedCourses int                 `json:"completed_courses"`
	GPA              string              `json:"gpa"`
}, error) {
	var enrollments []models.Enrollment
	err := s.db.Preload("Course").
		Where("student_id = ? AND grade IS NOT NULL", studentID).
		Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	result := struct {
		Enrollments      []models.Enrollment `json:"enrollments"`
		TotalCredits     decimal.Decimal     `json:"total_credits"`
		CompletedCourses int                 `json:"completed_courses"`
		GPA              string              `json:"gpa"`
	}{
		Enrollments:      enrollments,
		CompletedCourses: len(enrollments),
	}

	var totalCredits decimal.Decimal
	var totalGradePoints decimal.Decimal

	for _, enrollment := range enrollments {
		if enrollment.Grade == nil {
			continue
		}
		credits := decimal.NewFromInt(int64(enrollment.Course.Credits))
		totalCredits = totalCredits.Add(credits)
		totalGradePoints = totalGradePoints.Add(enrollment.Grade.Mul(credits))
	}

	result.TotalCredits = totalCredits
	if totalCredits.IsPositive() {
		result.GPA = calculateGPA(totalCredits, totalGradePoints)
	} else {
		result.GPA = "0.00"
	}

	return &result, nil
}

// GetCourseGradeStats 实现获取课程成绩统计
func (s *EnrollmentServiceImpl) GetCourseGradeStats(courseID, teacherID uint) (*struct {
	TotalStudents int             `json:"total_students"`
	GradedCount   int             `json:"graded_count"`
	HighestGrade  decimal.Decimal `json:"highest_grade"`
	LowestGrade   decimal.Decimal `json:"lowest_grade"`
	AverageGrade  decimal.Decimal `json:"average_grade"`
	Distribution  struct {
		A int `json:"a"`
		B int `json:"b"`
		C int `json:"c"`
		D int `json:"d"`
		F int `json:"f"`
	} `json:"distribution"`
}, error) {
	// 验证课程是否属于该教师
	course, err := s.courseRepo.FindByID(context.Background(), courseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}
	if course.TeacherID != teacherID {
		return nil, ErrUnauthorized
	}

	var enrollments []models.Enrollment
	err = s.db.Where("course_id = ?", courseID).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	result := struct {
		TotalStudents int             `json:"total_students"`
		GradedCount   int             `json:"graded_count"`
		HighestGrade  decimal.Decimal `json:"highest_grade"`
		LowestGrade   decimal.Decimal `json:"lowest_grade"`
		AverageGrade  decimal.Decimal `json:"average_grade"`
		Distribution  struct {
			A int `json:"a"`
			B int `json:"b"`
			C int `json:"c"`
			D int `json:"d"`
			F int `json:"f"`
		} `json:"distribution"`
	}{
		TotalStudents: len(enrollments),
		HighestGrade:  decimal.Zero,
		LowestGrade:   decimal.NewFromInt(100),
		AverageGrade:  decimal.Zero,
	}

	var totalGrade decimal.Decimal
	for _, enrollment := range enrollments {
		if enrollment.Grade == nil || enrollment.Grade.IsZero() {
			continue
		}

		result.GradedCount++
		totalGrade = totalGrade.Add(*enrollment.Grade)

		if enrollment.Grade.GreaterThan(result.HighestGrade) {
			result.HighestGrade = *enrollment.Grade
		}
		if enrollment.Grade.LessThan(result.LowestGrade) {
			result.LowestGrade = *enrollment.Grade
		}

		// 统计分数分布
		switch {
		case enrollment.Grade.GreaterThanOrEqual(decimal.NewFromInt(90)):
			result.Distribution.A++
		case enrollment.Grade.GreaterThanOrEqual(decimal.NewFromInt(80)):
			result.Distribution.B++
		case enrollment.Grade.GreaterThanOrEqual(decimal.NewFromInt(70)):
			result.Distribution.C++
		case enrollment.Grade.GreaterThanOrEqual(decimal.NewFromInt(60)):
			result.Distribution.D++
		default:
			result.Distribution.F++
		}
	}

	if result.GradedCount > 0 {
		result.AverageGrade = totalGrade.Div(decimal.NewFromInt(int64(result.GradedCount))).Round(2)
	}

	if result.GradedCount == 0 {
		result.HighestGrade = decimal.Zero
		result.LowestGrade = decimal.Zero
	}

	return &result, nil
}

func calculateGPA(totalCredits, totalGradePoints decimal.Decimal) string {
	if totalCredits.IsZero() {
		return "0.00"
	}
	gpa := totalGradePoints.Div(totalCredits).Round(2)
	return fmt.Sprintf("%.2f", gpa)
}
