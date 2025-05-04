package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/fynn404/go-vue-admin/internal/model"
)

// GradeHistoryRepository 定义成绩历史记录数据访问接口
type GradeHistoryRepository interface {
	Create(ctx context.Context, history *models.GradeHistory) error
	ListByGrade(ctx context.Context, gradeID uint, offset, limit int) ([]*models.GradeHistory, int64, error)
	ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.GradeHistory, int64, error)
	ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.GradeHistory, int64, error)
	ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.GradeHistory, int64, error)
	ListByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.GradeHistory, int64, error)
}

// gradeHistoryRepository 实现 GradeHistoryRepository 接口
type gradeHistoryRepository struct {
	db *gorm.DB
}

// NewGradeHistoryRepository 创建成绩历史记录 Repository 实例
func NewGradeHistoryRepository(db *gorm.DB) GradeHistoryRepository {
	return &gradeHistoryRepository{db: db}
}

// Create 创建成绩历史记录
func (r *gradeHistoryRepository) Create(ctx context.Context, history *models.GradeHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// ListByGrade 获取指定成绩的历史记录
func (r *gradeHistoryRepository) ListByGrade(ctx context.Context, gradeID uint, offset, limit int) ([]*models.GradeHistory, int64, error) {
	var histories []*models.GradeHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.GradeHistory{}).Where("grade_id = ?", gradeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("grade_id = ?", gradeID).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Order("operated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// ListByStudent 获取学生的成绩修改历史
func (r *gradeHistoryRepository) ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.GradeHistory, int64, error) {
	var histories []*models.GradeHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.GradeHistory{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Order("operated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// ListByCourse 获取课程的成绩修改历史
func (r *gradeHistoryRepository) ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.GradeHistory, int64, error) {
	var histories []*models.GradeHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.GradeHistory{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Order("operated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// ListByTeacher 获取教师的成绩修改历史
func (r *gradeHistoryRepository) ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.GradeHistory, int64, error) {
	var histories []*models.GradeHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.GradeHistory{}).Where("teacher_id = ?", teacherID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Order("operated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// ListByDateRange 获取指定日期范围内的成绩修改历史
func (r *gradeHistoryRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.GradeHistory, int64, error) {
	var histories []*models.GradeHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.GradeHistory{}).
		Where("operated_at BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("operated_at BETWEEN ? AND ?", startDate, endDate).
		Preload("Grade").
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Order("operated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
