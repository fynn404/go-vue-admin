package service

import (
	"context"
	"errors"

	models "github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
)

var (
	ErrCourseNotFound      = errors.New("course not found")
	ErrInvalidTeacher      = errors.New("invalid teacher")
	ErrInvalidCourseStatus = errors.New("invalid course status")
)

// CourseService 定义课程服务接口
type CourseService interface {
	// CreateCourse 创建新课程
	CreateCourse(ctx context.Context, course *models.Course) error
	// UpdateCourse 更新课程信息
	UpdateCourse(ctx context.Context, id uint, updates map[string]interface{}) error
	// DeleteCourse 删除课程
	DeleteCourse(ctx context.Context, id uint) error
	// GetCourseByID 根据ID获取课程信息
	GetCourseByID(ctx context.Context, id uint) (*models.Course, error)
	// ListCourses 获取课程列表
	ListCourses(ctx context.Context, page, pageSize int) ([]*models.Course, int64, error)
	// ListTeacherCourses 获取教师的课程列表
	ListTeacherCourses(ctx context.Context, teacherID uint, page, pageSize int) ([]*models.Course, int64, error)
	// ListAvailableCourses 获取可选课程列表
	ListAvailableCourses(ctx context.Context, page, pageSize int) ([]*models.Course, int64, error)
}

// courseService 实现 CourseService 接口
type courseService struct {
	courseRepo repository.CourseRepository
	userRepo   repository.UserRepository
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseRepo repository.CourseRepository, userRepo repository.UserRepository) CourseService {
	return &courseService{
		courseRepo: courseRepo,
		userRepo:   userRepo,
	}
}

// CreateCourse 实现创建课程
func (s *courseService) CreateCourse(ctx context.Context, course *models.Course) error {
	// 验证教师是否存在且角色正确
	teacher, err := s.userRepo.FindByID(ctx, course.TeacherID)
	if err != nil {
		return ErrInvalidTeacher
	}
	if teacher.Role != models.RoleTeacher {
		return ErrInvalidTeacher
	}

	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse 实现更新课程信息
func (s *courseService) UpdateCourse(ctx context.Context, id uint, updates map[string]interface{}) error {
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return ErrCourseNotFound
	}

	// 更新字段
	for key, value := range updates {
		switch key {
		case "name":
			course.Name = value.(string)
		case "description":
			course.Description = value.(string)
		case "credits":
			course.Credits = value.(float32)
		case "capacity":
			capacity := value.(int)
			if capacity < course.CurrentEnrolled {
				return models.ErrInvalidCapacity
			}
			course.Capacity = capacity
		case "status":
			status := models.CourseStatus(value.(string))
			if status != models.CourseStatusOpen && status != models.CourseStatusClosed {
				return ErrInvalidCourseStatus
			}
			course.Status = status
		}
	}

	return s.courseRepo.Update(ctx, course)
}

// DeleteCourse 实现删除课程
func (s *courseService) DeleteCourse(ctx context.Context, id uint) error {
	if _, err := s.courseRepo.FindByID(ctx, id); err != nil {
		return ErrCourseNotFound
	}
	return s.courseRepo.Delete(ctx, id)
}

// GetCourseByID 实现根据ID获取课程信息
func (s *courseService) GetCourseByID(ctx context.Context, id uint) (*models.Course, error) {
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

// ListCourses 实现获取课程列表
func (s *courseService) ListCourses(ctx context.Context, page, pageSize int) ([]*models.Course, int64, error) {
	offset := (page - 1) * pageSize
	return s.courseRepo.List(ctx, offset, pageSize)
}

// ListTeacherCourses 实现获取教师的课程列表
func (s *courseService) ListTeacherCourses(ctx context.Context, teacherID uint, page, pageSize int) ([]*models.Course, int64, error) {
	// 验证教师是否存在
	teacher, err := s.userRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, 0, ErrInvalidTeacher
	}
	if teacher.Role != models.RoleTeacher {
		return nil, 0, ErrInvalidTeacher
	}

	offset := (page - 1) * pageSize
	return s.courseRepo.ListByTeacher(ctx, teacherID, offset, pageSize)
}

// ListAvailableCourses 实现获取可选课程列表
func (s *courseService) ListAvailableCourses(ctx context.Context, page, pageSize int) ([]*models.Course, int64, error) {
	offset := (page - 1) * pageSize
	return s.courseRepo.ListAvailable(ctx, offset, pageSize)
}
