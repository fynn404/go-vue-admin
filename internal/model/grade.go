package model

import (
	"time"

	"github.com/shopspring/decimal"
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
	CourseID    uint            `gorm:"not null;index" json:"course_id"`
	Course      *Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	StudentID   uint            `gorm:"not null;index" json:"student_id"`
	Student     *User           `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	TeacherID   uint            `gorm:"not null;index" json:"teacher_id"`
	Teacher     *User           `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Score       decimal.Decimal `gorm:"type:decimal(5,2);not null" json:"score"`       // 分数，使用decimal类型，精确到小数点后2位
	GradePoint  decimal.Decimal `gorm:"type:decimal(3,2);not null" json:"grade_point"` // 绩点，使用decimal类型，精确到小数点后2位
	Status      GradeStatus     `gorm:"type:varchar(20);not null" json:"status"`       // 状态
	Comment     string          `gorm:"type:text" json:"comment"`                      // 评语
	PublishedAt *time.Time      `json:"published_at,omitempty"`                        // 发布时间
}

// CalculateGradePoint 计算绩点
func (g *Grade) CalculateGradePoint() {
	score := g.Score.InexactFloat64()
	switch {
	case score >= 90:
		g.GradePoint = decimal.NewFromFloat(4.0)
	case score >= 85:
		g.GradePoint = decimal.NewFromFloat(3.7)
	case score >= 80:
		g.GradePoint = decimal.NewFromFloat(3.3)
	case score >= 75:
		g.GradePoint = decimal.NewFromFloat(3.0)
	case score >= 70:
		g.GradePoint = decimal.NewFromFloat(2.7)
	case score >= 65:
		g.GradePoint = decimal.NewFromFloat(2.3)
	case score >= 60:
		g.GradePoint = decimal.NewFromFloat(2.0)
	default:
		g.GradePoint = decimal.NewFromFloat(0.0)
	}
}

// BeforeSave GORM 钩子，保存前计算绩点
func (g *Grade) BeforeSave(tx *gorm.DB) error {
	g.CalculateGradePoint()
	return nil
}

// TableName - Set table name for GORM
func (g *Grade) TableName() string {
	return "grades_tab"
}
