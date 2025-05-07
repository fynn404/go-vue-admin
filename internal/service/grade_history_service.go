package service

import (
	"time"

	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// GradeHistoryService 处理成绩历史记录相关的业务逻辑接口
type GradeHistoryService interface {
	// CreateHistory 创建成绩历史记录
	CreateHistory(gradeID, courseID, studentID, teacherID uint, changeType model.GradeChangeType, oldScore, newScore *decimal.Decimal, oldComment, newComment, reason string) (*model.GradeHistory, error)
	// GetHistoryByGradeID 获取指定成绩的历史记录
	GetHistoryByGradeID(gradeID uint) ([]model.GradeHistory, error)
	// GetHistoryByStudentID 获取指定学生的成绩历史记录
	GetHistoryByStudentID(studentID uint) ([]model.GradeHistory, error)
	// GetHistoryByCourseID 获取指定课程的成绩历史记录
	GetHistoryByCourseID(courseID uint) ([]model.GradeHistory, error)
	// GetHistoryByTeacherID 获取指定教师的操作历史记录
	GetHistoryByTeacherID(teacherID uint) ([]model.GradeHistory, error)
	// GetHistoryDetail 获取历史记录详情
	GetHistoryDetail(historyID uint) (*model.GradeHistory, error)
	// GetHistoryByDateRange 获取指定日期范围内的历史记录
	GetHistoryByDateRange(startDate, endDate time.Time) ([]model.GradeHistory, error)
}

// GradeHistoryServiceImpl 实现 GradeHistoryService 接口
type GradeHistoryServiceImpl struct {
	db *gorm.DB
}

// NewGradeHistoryService 创建成绩历史记录服务实例
func NewGradeHistoryService() GradeHistoryService {
	return &GradeHistoryServiceImpl{
		db: config.DB,
	}
}

// CreateHistory 创建成绩历史记录
func (s *GradeHistoryServiceImpl) CreateHistory(gradeID, courseID, studentID, teacherID uint, changeType model.GradeChangeType, oldScore, newScore *decimal.Decimal, oldComment, newComment, reason string) (*model.GradeHistory, error) {
	history := &model.GradeHistory{
		GradeID:    gradeID,
		CourseID:   courseID,
		StudentID:  studentID,
		TeacherID:  teacherID,
		ChangeType: changeType,
		OldScore:   oldScore,
		NewScore:   newScore,
		OldComment: oldComment,
		NewComment: newComment,
		Reason:     reason,
		OperatedAt: time.Now(),
	}

	if err := s.db.Create(history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// GetHistoryByGradeID 获取指定成绩的历史记录
func (s *GradeHistoryServiceImpl) GetHistoryByGradeID(gradeID uint) ([]model.GradeHistory, error) {
	var histories []model.GradeHistory
	err := s.db.Where("grade_id = ?", gradeID).
		Order("operated_at DESC").
		Find(&histories).Error
	return histories, err
}

// GetHistoryByStudentID 获取指定学生的成绩历史记录
func (s *GradeHistoryServiceImpl) GetHistoryByStudentID(studentID uint) ([]model.GradeHistory, error) {
	var histories []model.GradeHistory
	err := s.db.Where("student_id = ?", studentID).
		Order("operated_at DESC").
		Preload("Grade").
		Preload("Course").
		Find(&histories).Error
	return histories, err
}

// GetHistoryByCourseID 获取指定课程的成绩历史记录
func (s *GradeHistoryServiceImpl) GetHistoryByCourseID(courseID uint) ([]model.GradeHistory, error) {
	var histories []model.GradeHistory
	err := s.db.Where("course_id = ?", courseID).
		Order("operated_at DESC").
		Preload("Grade").
		Preload("Student").
		Find(&histories).Error
	return histories, err
}

// GetHistoryByTeacherID 获取指定教师的操作历史记录
func (s *GradeHistoryServiceImpl) GetHistoryByTeacherID(teacherID uint) ([]model.GradeHistory, error) {
	var histories []model.GradeHistory
	err := s.db.Where("teacher_id = ?", teacherID).
		Order("operated_at DESC").
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Find(&histories).Error
	return histories, err
}

// GetHistoryDetail 获取历史记录详情
func (s *GradeHistoryServiceImpl) GetHistoryDetail(historyID uint) (*model.GradeHistory, error) {
	var history model.GradeHistory
	err := s.db.Where("id = ?", historyID).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetHistoryByDateRange 获取指定日期范围内的历史记录
func (s *GradeHistoryServiceImpl) GetHistoryByDateRange(startDate, endDate time.Time) ([]model.GradeHistory, error) {
	var histories []model.GradeHistory
	err := s.db.Where("operated_at BETWEEN ? AND ?", startDate, endDate).
		Order("operated_at DESC").
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Find(&histories).Error
	return histories, err
}
