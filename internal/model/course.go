package model

import (
	"time"

	"gorm.io/gorm"
)

type CourseStatus string

const (
	CourseStatusOpen     CourseStatus = "open"
	CourseStatusClosed   CourseStatus = "closed"
	CourseStatusArchived CourseStatus = "archived"
)

// Course 课程模型
type Course struct {
	gorm.Model
	Name            string       `gorm:"not null" json:"name"`
	Description     string       `json:"description"`
	TeacherID       uint         `gorm:"not null" json:"teacher_id"`
	Teacher         User         `json:"teacher,omitempty"`
	Credits         int          `gorm:"not null" json:"credits"`
	CurrentEnrolled int          `gorm:"not null" json:"current_enrolled"`
	Capacity        int          `gorm:"not null" json:"capacity"`
	Status          CourseStatus `gorm:"not null" json:"status"`
	Students        []User       `gorm:"many2many:enrollments;" json:"students,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

// TableName - Set table name for GORM
func (Course) TableName() string {
	return "courses_tab"
}

// BeforeCreate - GORM hook to validate course before creation
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.Capacity < len(c.Students) {
		return ErrInvalidCapacity
	}
	return nil
}

// IsAvailable 检查课程是否可选
func (c *Course) IsAvailable() bool {
	return c.Status == CourseStatusOpen
}

// Enroll - Enroll a student in the course
func (c *Course) Enroll(student User) error {
	if !c.IsAvailable() {
		return ErrCourseUnavailable
	}
	c.Students = append(c.Students, student)
	return nil
}

// Unenroll - Unenroll a student from the course
func (c *Course) Unenroll(student User) error {
	if len(c.Students) <= 0 {
		return ErrNoEnrollment
	}
	for i, s := range c.Students {
		if s.ID == student.ID {
			c.Students = append(c.Students[:i], c.Students[i+1:]...)
			c.CurrentEnrolled--
			return nil
		}
	}
	return ErrStudentNotFound
}
