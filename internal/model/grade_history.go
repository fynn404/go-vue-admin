package models

import (
	"time"

	"gorm.io/gorm"
)

// GradeHistory 记录成绩变更历史
type GradeHistory struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	EnrollmentID uint           `json:"enrollment_id" gorm:"not null"`
	Enrollment   Enrollment     `json:"enrollment" gorm:"foreignKey:EnrollmentID"`
	TeacherID    uint           `json:"teacher_id" gorm:"not null"`
	Teacher      User           `json:"teacher" gorm:"foreignKey:TeacherID"`
	Grade        float32        `json:"grade" gorm:"not null"`
	Comment      string         `json:"comment"`
	Timestamp    time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName - Set table name for GORM
func (GradeHistory) TableName() string {
	return "grade_histories"
}
