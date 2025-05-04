package repository

import (
	"context"

	"gorm.io/gorm"

	models "github.com/fynn404/go-vue-admin/internal/model"
)

// EnrollmentRepository 定义选课数据访问接口
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	Update(ctx context.Context, enrollment *models.Enrollment) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Enrollment, error)
	FindByStudentAndCourse(ctx context.Context, studentID, courseID uint) (*models.Enrollment, error)
	ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.Enrollment, int64, error)
	ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.Enrollment, int64, error)
}

// enrollmentRepository 实现 EnrollmentRepository 接口
type enrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository 创建选课 Repository 实例
func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

// Create 创建选课记录
func (r *enrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

// Update 更新选课记录
func (r *enrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}

// Delete 删除选课记录
func (r *enrollmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Enrollment{}, id).Error
}

// FindByID 根据ID查找选课记录
func (r *enrollmentRepository) FindByID(ctx context.Context, id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).Preload("Student").Preload("Course").First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// FindByStudentAndCourse 根据学生ID和课程ID查找选课记录
func (r *enrollmentRepository) FindByStudentAndCourse(ctx context.Context, studentID, courseID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Course").
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		First(&enrollment).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// ListByStudent 获取学生的选课列表
func (r *enrollmentRepository) ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.Enrollment, int64, error) {
	var enrollments []*models.Enrollment
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Course").
		Where("student_id = ?", studentID).
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

// ListByCourse 获取课程的选课列表
func (r *enrollmentRepository) ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.Enrollment, int64, error) {
	var enrollments []*models.Enrollment
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Course").
		Where("course_id = ?", courseID).
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}
