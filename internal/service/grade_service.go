package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrGradeNotFound     = errors.New("grade not found")
	ErrNotAuthorized     = errors.New("not authorized")
	ErrInvalidStudentID  = errors.New("invalid student ID")
	ErrInvalidCourseID   = errors.New("invalid course ID")
	ErrStudentNotFound   = errors.New("student not found")
	ErrGradeAlreadyExist = errors.New("grade already exists for this student in this course")
	ErrInvalidScore      = errors.New("score must be between 0 and 100")
)

// GradeService 处理成绩相关的业务逻辑接口
type GradeService interface {
	// CreateGradeWithValidation creates a grade with full validation
	CreateGradeWithValidation(courseIDStr, studentIDStr string, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error)
	// UpdateGradeWithValidation updates a grade with full validation
	UpdateGradeWithValidation(gradeIDStr string, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error)
	// CreateGrade 创建成绩
	CreateGrade(courseID, studentID, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error)
	// UpdateGrade 更新成绩
	UpdateGrade(gradeID, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error)
	// PublishGrade 发布成绩
	PublishGrade(gradeID, teacherID uint) (*model.Grade, error)
	// GetStudentGrades 获取学生成绩
	GetStudentGrades(studentID uint) ([]model.Grade, error)
	// GetCourseGrades 获取课程成绩
	GetCourseGrades(courseID, teacherID uint) ([]model.Grade, error)
	// GetGradeHistory 获取成绩历史
	GetGradeHistory(gradeID, userID uint, role string) ([]model.GradeHistory, error)
}

// GradeServiceImpl 实现 GradeService 接口
type GradeServiceImpl struct {
	db *gorm.DB
}

// NewGradeService 创建成绩服务实例
func NewGradeService() GradeService {
	return &GradeServiceImpl{
		db: config.DB,
	}
}

// validateScore checks if the score is within valid range
func (s *GradeServiceImpl) validateScore(score decimal.Decimal) error {
	if score.LessThan(decimal.Zero) || score.GreaterThan(decimal.NewFromInt(100)) {
		return ErrInvalidScore
	}
	return nil
}

// ParseAndValidateID parses and validates a string ID
func (s *GradeServiceImpl) ParseAndValidateID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return uint(parsed), nil
}

// CreateGradeWithValidation creates a grade with full validation
func (s *GradeServiceImpl) CreateGradeWithValidation(courseIDStr, studentIDStr string, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error) {
	// Validate score
	if err := s.validateScore(score); err != nil {
		return nil, err
	}

	// Parse and validate IDs
	courseID, err := s.ParseAndValidateID(courseIDStr)
	if err != nil {
		return nil, ErrInvalidCourseID
	}

	studentID, err := s.ParseAndValidateID(studentIDStr)
	if err != nil {
		return nil, ErrInvalidStudentID
	}

	return s.CreateGrade(courseID, studentID, teacherID, score, comment)
}

// UpdateGradeWithValidation updates a grade with full validation
func (s *GradeServiceImpl) UpdateGradeWithValidation(gradeIDStr string, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error) {
	// Validate score
	if err := s.validateScore(score); err != nil {
		return nil, err
	}

	// Parse and validate grade ID
	gradeID, err := s.ParseAndValidateID(gradeIDStr)
	if err != nil {
		return nil, ErrGradeNotFound
	}

	return s.UpdateGrade(gradeID, teacherID, score, comment)
}

// CreateGrade 创建成绩
func (s *GradeServiceImpl) CreateGrade(courseID, studentID, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error) {
	// 检查课程是否存在
	var course model.Course
	if err := s.db.First(&course, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	// 检查教师权限
	if course.TeacherID != teacherID {
		return nil, ErrNotAuthorized
	}

	// 检查学生是否存在
	var student model.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	// 检查是否已存在成绩
	var existingGrade model.Grade
	err := s.db.Where("course_id = ? AND student_id = ?", courseID, studentID).First(&existingGrade).Error
	if err == nil {
		return nil, ErrGradeAlreadyExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 创建成绩记录
	grade := &model.Grade{
		CourseID:  courseID,
		StudentID: studentID,
		TeacherID: teacherID,
		Score:     score,
		Comment:   comment,
		Status:    model.GradeStatusDraft,
	}

	tx := s.db.Begin()

	if err := tx.Create(grade).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建历史记录
	history := model.GradeHistory{
		GradeID:    grade.ID,
		CourseID:   grade.CourseID,
		StudentID:  grade.StudentID,
		TeacherID:  grade.TeacherID,
		ChangeType: model.GradeChangeTypeCreate,
		NewScore:   &grade.Score,
		NewComment: grade.Comment,
		OperatedAt: time.Now(),
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return grade, nil
}

// UpdateGrade 更新成绩
func (s *GradeServiceImpl) UpdateGrade(gradeID, teacherID uint, score decimal.Decimal, comment string) (*model.Grade, error) {
	var grade model.Grade
	if err := s.db.First(&grade, gradeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}

	// 检查教师权限
	if grade.TeacherID != teacherID {
		return nil, ErrNotAuthorized
	}

	tx := s.db.Begin()

	// 创建历史记录
	history := model.GradeHistory{
		GradeID:    grade.ID,
		CourseID:   grade.CourseID,
		StudentID:  grade.StudentID,
		TeacherID:  teacherID,
		ChangeType: model.GradeChangeTypeUpdate,
		OldScore:   &grade.Score,
		NewScore:   &score,
		OldComment: grade.Comment,
		NewComment: comment,
		OperatedAt: time.Now(),
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新成绩
	grade.Score = score
	grade.Comment = comment

	if err := tx.Save(&grade).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &grade, nil
}

// PublishGrade 发布成绩
func (s *GradeServiceImpl) PublishGrade(gradeID, teacherID uint) (*model.Grade, error) {
	var grade model.Grade
	if err := s.db.First(&grade, gradeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}

	// 检查教师权限
	if grade.TeacherID != teacherID {
		return nil, ErrNotAuthorized
	}

	// 更新状态为已发布
	grade.Status = model.GradeStatusPublished

	if err := s.db.Save(&grade).Error; err != nil {
		return nil, err
	}

	return &grade, nil
}

// GetStudentGrades 获取学生成绩
func (s *GradeServiceImpl) GetStudentGrades(studentID uint) ([]model.Grade, error) {
	var grades []model.Grade
	if err := s.db.Where("student_id = ?", studentID).Find(&grades).Error; err != nil {
		return nil, err
	}
	return grades, nil
}

// GetCourseGrades 获取课程成绩
func (s *GradeServiceImpl) GetCourseGrades(courseID, teacherID uint) ([]model.Grade, error) {
	// 检查教师权限
	var course model.Course
	if err := s.db.First(&course, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	if course.TeacherID != teacherID {
		return nil, ErrNotAuthorized
	}

	var grades []model.Grade
	if err := s.db.Where("course_id = ?", courseID).Find(&grades).Error; err != nil {
		return nil, err
	}

	return grades, nil
}

// GetGradeHistory 获取成绩历史
func (s *GradeServiceImpl) GetGradeHistory(gradeID, userID uint, role string) ([]model.GradeHistory, error) {
	var grade model.Grade
	if err := s.db.First(&grade, gradeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}

	// 检查权限
	switch role {
	case "admin":
		// 管理员可以查看所有历史记录
	case "teacher":
		if grade.TeacherID != userID {
			return nil, ErrNotAuthorized
		}
	case "student":
		if grade.StudentID != userID {
			return nil, ErrNotAuthorized
		}
	default:
		return nil, ErrNotAuthorized
	}

	var histories []model.GradeHistory
	if err := s.db.Where("grade_id = ?", gradeID).Order("operated_at DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}
