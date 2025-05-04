package models

import (
	"time"

	"gorm.io/gorm"
)

// GradeStatus 成绩状态
type GradeStatus string

const (
	GradeStatusDraft     GradeStatus = "draft"     // 草稿
	GradeStatusPublished GradeStatus = "published" // 已发布
)

// Grade 成绩模型
type Grade struct {
	gorm.Model
	CourseID    uint        `gorm:"not null;index" json:"course_id"`
	Course      *Course     `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	StudentID   uint        `gorm:"not null;index" json:"student_id"`
	Student     *User       `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	TeacherID   uint        `gorm:"not null;index" json:"teacher_id"`
	Teacher     *User       `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Score       float32     `gorm:"not null" json:"score"`                   // 分数
	GradePoint  float32     `gorm:"not null" json:"grade_point"`             // 绩点
	Status      GradeStatus `gorm:"type:varchar(20);not null" json:"status"` // 状态
	Comment     string      `gorm:"type:text" json:"comment"`                // 评语
	PublishedAt *time.Time  `json:"published_at,omitempty"`                  // 发布时间
}

// CalculateGradePoint 计算绩点
func (g *Grade) CalculateGradePoint() {
	switch {
	case g.Score >= 90:
		g.GradePoint = 4.0
	case g.Score >= 85:
		g.GradePoint = 3.7
	case g.Score >= 80:
		g.GradePoint = 3.3
	case g.Score >= 75:
		g.GradePoint = 3.0
	case g.Score >= 70:
		g.GradePoint = 2.7
	case g.Score >= 65:
		g.GradePoint = 2.3
	case g.Score >= 60:
		g.GradePoint = 2.0
	default:
		g.GradePoint = 0.0
	}
}

// BeforeSave GORM 钩子，保存前计算绩点
func (g *Grade) BeforeSave(tx *gorm.DB) error {
	g.CalculateGradePoint()
	return nil
}
