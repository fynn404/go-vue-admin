package service

import (
	"context"
	"errors"

	models "github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
)

var (
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrInvalidStudent     = errors.New("invalid student")
	ErrInvalidGrade       = errors.New("invalid grade")
)

// EnrollmentService 定义选课服务接口
type EnrollmentService interface {
	// EnrollCourse 学生选课
	EnrollCourse(ctx context.Context, studentID, courseID uint) error
	// DropCourse 退课
	DropCourse(ctx context.Context, enrollmentID uint) error
	// UpdateGrade 更新成绩
	UpdateGrade(ctx context.Context, enrollmentID uint, grade float32) error
	// GetEnrollment 获取选课记录
	GetEnrollment(ctx context.Context, enrollmentID uint) (*models.Enrollment, error)
	// ListStudentEnrollments 获取学生的选课列表
	ListStudentEnrollments(ctx context.Context, studentID uint, page, pageSize int) ([]*models.Enrollment, int64, error)
	// ListCourseEnrollments 获取课程的选课列表
	ListCourseEnrollments(ctx context.Context, courseID uint, page, pageSize int) ([]*models.Enrollment, int64, error)
}

// enrollmentService 实现 EnrollmentService 接口
type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
}

// NewEnrollmentService 创建选课服务实例
func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	userRepo repository.UserRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
	}
}

// EnrollCourse 实现学生选课
func (s *enrollmentService) EnrollCourse(ctx context.Context, studentID, courseID uint) error {
	// 验证学生身份
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return ErrInvalidStudent
	}
	if student.Role != models.RoleStudent {
		return ErrInvalidStudent
	}

	// 验证课程是否存在并可选
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return models.ErrCourseUnavailable
	}
	if !course.IsAvailable() {
		return models.ErrCourseUnavailable
	}

	// 创建选课记录
	enrollment := &models.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    models.EnrollmentStatusActive,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return err
	}

	// 更新课程选课人数
	course.CurrentEnrolled++
	return s.courseRepo.Update(ctx, course)
}

// DropCourse 实现退课
func (s *enrollmentService) DropCourse(ctx context.Context, enrollmentID uint) error {
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

	if err := course.Unenroll(); err != nil {
		return err
	}

	return s.courseRepo.Update(ctx, course)
}

// UpdateGrade 实现更新成绩
func (s *enrollmentService) UpdateGrade(ctx context.Context, enrollmentID uint, grade float32) error {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return ErrEnrollmentNotFound
	}

	if err := enrollment.UpdateGrade(grade); err != nil {
		return ErrInvalidGrade
	}

	return s.enrollmentRepo.Update(ctx, enrollment)
}

// GetEnrollment 实现获取选课记录
func (s *enrollmentService) GetEnrollment(ctx context.Context, enrollmentID uint) (*models.Enrollment, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}
	return enrollment, nil
}

// ListStudentEnrollments 实现获取学生的选课列表
func (s *enrollmentService) ListStudentEnrollments(ctx context.Context, studentID uint, page, pageSize int) ([]*models.Enrollment, int64, error) {
	// 验证学生是否存在
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, 0, ErrInvalidStudent
	}
	if student.Role != models.RoleStudent {
		return nil, 0, ErrInvalidStudent
	}

	offset := (page - 1) * pageSize
	return s.enrollmentRepo.ListByStudent(ctx, studentID, offset, pageSize)
}

// ListCourseEnrollments 实现获取课程的选课列表
func (s *enrollmentService) ListCourseEnrollments(ctx context.Context, courseID uint, page, pageSize int) ([]*models.Enrollment, int64, error) {
	// 验证课程是否存在
	if _, err := s.courseRepo.FindByID(ctx, courseID); err != nil {
		return nil, 0, models.ErrCourseUnavailable
	}

	offset := (page - 1) * pageSize
	return s.enrollmentRepo.ListByCourse(ctx, courseID, offset, pageSize)
}
