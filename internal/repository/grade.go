package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	models "github.com/fynn404/go-vue-admin/internal/model"
)

// GradeRepository 定义成绩数据访问接口
type GradeRepository interface {
	// 基础操作
	Create(ctx context.Context, grade *models.Grade) error
	Update(ctx context.Context, grade *models.Grade) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Grade, error)

	// 查询操作
	List(ctx context.Context, offset, limit int) ([]*models.Grade, int64, error)
	ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.Grade, int64, error)
	ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.Grade, int64, error)
	ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.Grade, int64, error)

	// 成绩统计
	GetStudentGPA(ctx context.Context, studentID uint) (float32, error)
	GetCourseStats(ctx context.Context, courseID uint) (min, max, avg float32, err error)

	// 批量操作
	BatchCreate(ctx context.Context, grades []*models.Grade) error
	BatchPublish(ctx context.Context, ids []uint) error
}

// gradeRepository 实现 GradeRepository 接口
type gradeRepository struct {
	db *gorm.DB
}

// NewGradeRepository 创建成绩 Repository 实例
func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db: db}
}

// Create 创建成绩记录
func (r *gradeRepository) Create(ctx context.Context, grade *models.Grade) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建成绩记录
		if err := tx.Create(grade).Error; err != nil {
			return err
		}

		// 创建历史记录
		history := &models.GradeHistory{
			GradeID:    grade.ID,
			CourseID:   grade.CourseID,
			StudentID:  grade.StudentID,
			TeacherID:  grade.TeacherID,
			ChangeType: models.GradeChangeTypeCreate,
			NewScore:   &grade.Score,
			NewComment: grade.Comment,
			OperatedAt: time.Now(),
		}

		// 保存历史记录
		return tx.Create(history).Error
	})
}

// Update 更新成绩信息
func (r *gradeRepository) Update(ctx context.Context, grade *models.Grade) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取原成绩记录
		var oldGrade models.Grade
		if err := tx.First(&oldGrade, grade.ID).Error; err != nil {
			return err
		}

		// 创建历史记录
		history := &models.GradeHistory{
			GradeID:    grade.ID,
			CourseID:   grade.CourseID,
			StudentID:  grade.StudentID,
			TeacherID:  grade.TeacherID,
			ChangeType: models.GradeChangeTypeUpdate,
			OldScore:   &oldGrade.Score,
			NewScore:   &grade.Score,
			OldComment: oldGrade.Comment,
			NewComment: grade.Comment,
			OperatedAt: time.Now(),
		}

		// 保存历史记录
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		// 更新成绩
		return tx.Save(grade).Error
	})
}

// Delete 删除成绩记录
func (r *gradeRepository) Delete(ctx context.Context, id uint) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取原成绩记录
		var grade models.Grade
		if err := tx.First(&grade, id).Error; err != nil {
			return err
		}

		// 创建历史记录
		history := &models.GradeHistory{
			GradeID:    grade.ID,
			CourseID:   grade.CourseID,
			StudentID:  grade.StudentID,
			TeacherID:  grade.TeacherID,
			ChangeType: models.GradeChangeTypeDelete,
			OldScore:   &grade.Score,
			OldComment: grade.Comment,
			OperatedAt: time.Now(),
		}

		// 保存历史记录
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		// 删除成绩
		return tx.Delete(&grade).Error
	})
}

// FindByID 根据ID查找成绩
func (r *gradeRepository) FindByID(ctx context.Context, id uint) (*models.Grade, error) {
	var grade models.Grade
	if err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		First(&grade, id).Error; err != nil {
		return nil, err
	}
	return &grade, nil
}

// List 获取成绩列表
func (r *gradeRepository) List(ctx context.Context, offset, limit int) ([]*models.Grade, int64, error) {
	var grades []*models.Grade
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Grade{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Offset(offset).
		Limit(limit).
		Find(&grades).Error; err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// ListByCourse 获取课程的成绩列表
func (r *gradeRepository) ListByCourse(ctx context.Context, courseID uint, offset, limit int) ([]*models.Grade, int64, error) {
	var grades []*models.Grade
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Grade{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Offset(offset).
		Limit(limit).
		Find(&grades).Error; err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// ListByStudent 获取学生的成绩列表
func (r *gradeRepository) ListByStudent(ctx context.Context, studentID uint, offset, limit int) ([]*models.Grade, int64, error) {
	var grades []*models.Grade
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Grade{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Offset(offset).
		Limit(limit).
		Find(&grades).Error; err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// ListByTeacher 获取教师评定的成绩列表
func (r *gradeRepository) ListByTeacher(ctx context.Context, teacherID uint, offset, limit int) ([]*models.Grade, int64, error) {
	var grades []*models.Grade
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Grade{}).Where("teacher_id = ?", teacherID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Preload("Course").
		Preload("Student").
		Preload("Teacher").
		Offset(offset).
		Limit(limit).
		Find(&grades).Error; err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// GetStudentGPA 计算学生的平均绩点
func (r *gradeRepository) GetStudentGPA(ctx context.Context, studentID uint) (float32, error) {
	type Result struct {
		TotalCredits    float32
		TotalGradePoint float32
	}
	var result Result

	err := r.db.WithContext(ctx).
		Table("grades").
		Select("SUM(courses.credits) as total_credits, SUM(courses.credits * grades.grade_point) as total_grade_point").
		Joins("JOIN courses ON grades.course_id = courses.id").
		Where("grades.student_id = ? AND grades.status = ?", studentID, models.GradeStatusPublished).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	if result.TotalCredits == 0 {
		return 0, nil
	}

	return result.TotalGradePoint / result.TotalCredits, nil
}

// GetCourseStats 获取课程成绩统计信息
func (r *gradeRepository) GetCourseStats(ctx context.Context, courseID uint) (min, max, avg float32, err error) {
	type Stats struct {
		Min float32
		Max float32
		Avg float32
	}
	var stats Stats

	err = r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Select("MIN(score) as min, MAX(score) as max, AVG(score) as avg").
		Where("course_id = ? AND status = ?", courseID, models.GradeStatusPublished).
		Scan(&stats).Error

	if err != nil {
		return 0, 0, 0, err
	}

	return stats.Min, stats.Max, stats.Avg, nil
}

// BatchCreate 批量创建成绩记录
func (r *gradeRepository) BatchCreate(ctx context.Context, grades []*models.Grade) error {
	return r.db.WithContext(ctx).Create(&grades).Error
}

// BatchPublish 批量发布成绩
func (r *gradeRepository) BatchPublish(ctx context.Context, ids []uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":       models.GradeStatusPublished,
			"published_at": now,
		}).Error
}
