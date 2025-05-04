package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type EnrollmentStatus string

const (
	EnrollmentStatusActive   EnrollmentStatus = "active"
	EnrollmentStatusDropped  EnrollmentStatus = "dropped"
	EnrollmentStatusComplete EnrollmentStatus = "complete"
)

type Enrollment struct {
	ID        uint             `json:"id" gorm:"primarykey"`
	StudentID uint             `json:"student_id" gorm:"not null"`
	Student   User             `json:"student" gorm:"foreignKey:StudentID"`
	CourseID  uint             `json:"course_id" gorm:"not null"`
	Course    Course           `json:"course" gorm:"foreignKey:CourseID"`
	Status    EnrollmentStatus `json:"status" gorm:"type:varchar(10);not null;default:'active'"`
	Grade     *float32         `json:"grade"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`
}

// TableName - Set table name for GORM
func (Enrollment) TableName() string {
	return "enrollments"
}

// BeforeCreate - GORM hook to validate enrollment before creation
func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	var course Course
	if err := tx.First(&course, e.CourseID).Error; err != nil {
		return err
	}

	if !course.IsAvailable() {
		return ErrCourseUnavailable
	}

	var count int64
	if err := tx.Model(&Enrollment{}).
		Where("student_id = ? AND course_id = ? AND status = ?",
			e.StudentID, e.CourseID, EnrollmentStatusActive).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrAlreadyEnrolled
	}

	return nil
}

// UpdateGrade - Update the grade for an enrollment
func (e *Enrollment) UpdateGrade(grade float32) error {
	if grade < 0 || grade > 100 {
		return errors.New("grade must be between 0 and 100")
	}
	e.Grade = &grade
	return nil
}

// Drop - Drop the course
func (e *Enrollment) Drop() error {
	if e.Status != EnrollmentStatusActive {
		return errors.New("can only drop active enrollments")
	}
	e.Status = EnrollmentStatusDropped
	return nil
}
