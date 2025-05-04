package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	models "github.com/fynn404/go-vue-admin/internal/model"
)

// CourseRepository 定义课程数据访问接口
type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Course, error)
	List(ctx context.Context, offset, limit int) ([]*models.Course, int64, error)
	ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.Course, int64, error)
	ListAvailable(ctx context.Context, offset, limit int) ([]*models.Course, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*models.Course, int64, error)
	ListByCreditsRange(ctx context.Context, minCredits, maxCredits float32, offset, limit int) ([]*models.Course, int64, error)
	ListByStatus(ctx context.Context, status models.CourseStatus, offset, limit int) ([]*models.Course, int64, error)
	BatchDelete(ctx context.Context, ids []uint) error
}

// courseRepository 实现 CourseRepository 接口
type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository 创建课程 Repository 实例
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

// Create 创建课程
func (r *courseRepository) Create(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

// Update 更新课程信息
func (r *courseRepository) Update(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

// Delete 删除课程
func (r *courseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Course{}, id).Error
}

// FindByID 根据ID查找课程
func (r *courseRepository) FindByID(ctx context.Context, id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Preload("Teacher").First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

// List 获取课程列表
func (r *courseRepository) List(ctx context.Context, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Teacher").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// ListByTeacher 获取教师的课程列表
func (r *courseRepository) ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Course{}).Where("teacher_id = ?", teacherID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Teacher").Where("teacher_id = ?", teacherID).Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// ListAvailable 获取可选课程列表
func (r *courseRepository) ListAvailable(ctx context.Context, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	query := r.db.WithContext(ctx).Where("status = ? AND current_enrolled < capacity", models.CourseStatusOpen)

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Teacher").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// Search 搜索课程
func (r *courseRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	// 构建模糊搜索条件
	keyword = "%" + strings.ToLower(keyword) + "%"
	query := r.db.WithContext(ctx).Where(
		"LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
		keyword, keyword,
	)

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Teacher").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// ListByCreditsRange 按学分范围查询课程
func (r *courseRepository) ListByCreditsRange(ctx context.Context, minCredits, maxCredits float32, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	query := r.db.WithContext(ctx).Where("credits BETWEEN ? AND ?", minCredits, maxCredits)

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Teacher").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// ListByStatus 按状态查询课程
func (r *courseRepository) ListByStatus(ctx context.Context, status models.CourseStatus, offset, limit int) ([]*models.Course, int64, error) {
	var courses []*models.Course
	var total int64

	query := r.db.WithContext(ctx).Where("status = ?", status)

	if err := query.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Teacher").Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// BatchDelete 批量删除课程
func (r *courseRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&models.Course{}, ids).Error
}
