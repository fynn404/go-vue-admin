package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Enrollment 选课记录模型
type Enrollment struct {
	gorm.Model
	StudentID  uint             `gorm:"not null" json:"student_id"`
	Student    User             `json:"student,omitempty"`
	CourseID   uint             `gorm:"not null" json:"course_id"`
	Course     Course           `json:"course,omitempty"`
	Grade      *decimal.Decimal `gorm:"type:decimal(5,2)" json:"grade"`
	Status     EnrollmentStatus `gorm:"not null;default:'enrolled'" json:"status"`
	EnrolledAt time.Time        `json:"enrolled_at"`
	DroppedAt  *time.Time       `json:"dropped_at,omitempty"`
}

// TableName - Set table name for GORM
func (Enrollment) TableName() string {
	return "enrollment_tab"
}

// BeforeCreate - GORM hook
func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Enrollment{}).
		Where("student_id = ? AND course_id = ? AND status = ?",
			e.StudentID, e.CourseID, EnrollmentStatusEnrolled).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrDuplicateEnrollment
	}

	e.EnrolledAt = time.Now()
	return nil
}

// UpdateGrade - 更新成绩
func (e *Enrollment) UpdateGrade(grade decimal.Decimal) error {
	if e.Status != EnrollmentStatusEnrolled {
		return ErrInvalidEnrollmentStatus
	}
	e.Grade = &grade
	return nil
}

// Drop - 退课
func (e *Enrollment) Drop() error {
	if e.Status != EnrollmentStatusEnrolled {
		return ErrInvalidEnrollmentStatus
	}
	e.Status = EnrollmentStatusDropped
	now := time.Now()
	e.DroppedAt = &now
	return nil
}
