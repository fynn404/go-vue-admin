package models

import (
	"time"

	"gorm.io/gorm"
)

// GradeChangeType 成绩变更类型
type GradeChangeType string

const (
	GradeChangeTypeCreate GradeChangeType = "create" // 创建
	GradeChangeTypeUpdate GradeChangeType = "update" // 更新
	GradeChangeTypeDelete GradeChangeType = "delete" // 删除
)

// GradeHistory 成绩修改历史记录
type GradeHistory struct {
	gorm.Model
	GradeID    uint            `gorm:"not null;index" json:"grade_id"`               // 成绩ID
	CourseID   uint            `gorm:"not null;index" json:"course_id"`              // 课程ID
	StudentID  uint            `gorm:"not null;index" json:"student_id"`             // 学生ID
	TeacherID  uint            `gorm:"not null;index" json:"teacher_id"`             // 操作教师ID
	ChangeType GradeChangeType `gorm:"type:varchar(20);not null" json:"change_type"` // 变更类型
	OldScore   *float32        `json:"old_score,omitempty"`                          // 原分数
	NewScore   *float32        `json:"new_score,omitempty"`                          // 新分数
	OldComment string          `gorm:"type:text" json:"old_comment"`                 // 原评语
	NewComment string          `gorm:"type:text" json:"new_comment"`                 // 新评语
	Reason     string          `gorm:"type:text" json:"reason"`                      // 修改原因
	OperatedAt time.Time       `gorm:"not null" json:"operated_at"`                  // 操作时间

	// 关联
	Grade   *Grade  `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	Course  *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Student *User   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher *User   `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// TableName - Set table name for GORM
func (GradeHistory) TableName() string {
	return "grade_histories"
}
