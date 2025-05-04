package models

import (
	"time"

	"gorm.io/gorm"
)

type CourseStatus string

const (
	CourseStatusOpen   CourseStatus = "open"
	CourseStatusClosed CourseStatus = "closed"
)

type Course struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description"`
	TeacherID       uint           `json:"teacher_id" gorm:"not null"`
	Teacher         User           `json:"teacher" gorm:"foreignKey:TeacherID"`
	Credits         float32        `json:"credits" gorm:"not null"`
	Capacity        int            `json:"capacity" gorm:"not null"`
	CurrentEnrolled int            `json:"current_enrolled" gorm:"default:0"`
	Status          CourseStatus   `json:"status" gorm:"type:varchar(10);not null;default:'open'"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName - Set table name for GORM
func (Course) TableName() string {
	return "courses"
}

// BeforeCreate - GORM hook to validate course before creation
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.Capacity < c.CurrentEnrolled {
		return ErrInvalidCapacity
	}
	return nil
}

// IsAvailable - Check if course is available for enrollment
func (c *Course) IsAvailable() bool {
	return c.Status == CourseStatusOpen && c.CurrentEnrolled < c.Capacity
}

// Enroll - Enroll a student in the course
func (c *Course) Enroll() error {
	if !c.IsAvailable() {
		return ErrCourseUnavailable
	}
	c.CurrentEnrolled++
	return nil
}

// Unenroll - Unenroll a student from the course
func (c *Course) Unenroll() error {
	if c.CurrentEnrolled <= 0 {
		return ErrNoEnrollment
	}
	c.CurrentEnrolled--
	return nil
}
