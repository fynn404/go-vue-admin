package service

import (
	"context"
	"fmt"
	"math"

	"github.com/fynn404/go-vue-admin/internal/repository"
)

// GradeDistribution 成绩分布
type GradeDistribution struct {
	Range      string  `json:"range"`       // 分数段
	Count      int     `json:"count"`       // 人数
	Percent    float64 `json:"percent"`     // 百分比
	StartScore float64 `json:"start_score"` // 起始分数
	EndScore   float64 `json:"end_score"`   // 结束分数
}

// CourseGradeStats 课程成绩统计
type CourseGradeStats struct {
	CourseID      uint                `json:"course_id"`
	CourseName    string              `json:"course_name"`
	TotalStudents int                 `json:"total_students"`
	HighestScore  float32             `json:"highest_score"`
	LowestScore   float32             `json:"lowest_score"`
	AverageScore  float32             `json:"average_score"`
	MedianScore   float32             `json:"median_score"`
	PassRate      float64             `json:"pass_rate"`
	Distribution  []GradeDistribution `json:"distribution"`
}

// GradeStatsService 成绩统计服务
type GradeStatsService struct {
	gradeRepo repository.GradeRepository
}

// NewGradeStatsService 创建成绩统计服务实例
func NewGradeStatsService(gradeRepo repository.GradeRepository) *GradeStatsService {
	return &GradeStatsService{
		gradeRepo: gradeRepo,
	}
}

// GetCourseStats 获取课程成绩统计信息
func (s *GradeStatsService) GetCourseStats(ctx context.Context, courseID uint) (*CourseGradeStats, error) {
	// 获取课程所有成绩
	grades, total, err := s.gradeRepo.ListByCourse(ctx, courseID, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("获取课程成绩失败: %w", err)
	}

	if total == 0 || len(grades) == 0 {
		return nil, fmt.Errorf("未找到课程成绩数据")
	}

	// 初始化统计数据
	stats := &CourseGradeStats{
		CourseID:      courseID,
		CourseName:    grades[0].Course.Name,
		TotalStudents: int(total),
	}

	// 计算基础统计数据
	var sum float32
	var scores []float32
	passCount := 0

	for _, grade := range grades {
		scores = append(scores, grade.Score)
		sum += grade.Score
		if grade.Score >= 60 {
			passCount++
		}

		// 更新最高分和最低分
		if stats.HighestScore == 0 || grade.Score > stats.HighestScore {
			stats.HighestScore = grade.Score
		}
		if stats.LowestScore == 0 || grade.Score < stats.LowestScore {
			stats.LowestScore = grade.Score
		}
	}

	// 计算平均分
	stats.AverageScore = sum / float32(len(grades))

	// 计算及格率
	stats.PassRate = float64(passCount) / float64(len(grades)) * 100

	// 计算分数分布
	stats.Distribution = s.calculateDistribution(scores)

	return stats, nil
}

// calculateDistribution 计算成绩分布
func (s *GradeStatsService) calculateDistribution(scores []float32) []GradeDistribution {
	// 定义分数段
	ranges := []struct {
		start float64
		end   float64
		name  string
	}{
		{90, 100, "优秀"},
		{80, 89.99, "良好"},
		{70, 79.99, "中等"},
		{60, 69.99, "及格"},
		{0, 59.99, "不及格"},
	}

	// 初始化分布统计
	distribution := make([]GradeDistribution, len(ranges))
	totalCount := len(scores)

	// 统计每个分数段的人数
	for i, r := range ranges {
		count := 0
		for _, score := range scores {
			if float64(score) >= r.start && float64(score) <= r.end {
				count++
			}
		}
		percent := float64(count) / float64(totalCount) * 100

		distribution[i] = GradeDistribution{
			Range:      r.name,
			Count:      count,
			Percent:    math.Round(percent*100) / 100, // 保留两位小数
			StartScore: r.start,
			EndScore:   r.end,
		}
	}

	return distribution
}

// GetStudentGradeStats 获取学生成绩统计信息
func (s *GradeStatsService) GetStudentGradeStats(ctx context.Context, studentID uint) (map[string]interface{}, error) {
	// 获取学生所有成绩
	grades, total, err := s.gradeRepo.ListByStudent(ctx, studentID, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("获取学生成绩失败: %w", err)
	}

	if total == 0 || len(grades) == 0 {
		return nil, fmt.Errorf("未找到学生成绩数据")
	}

	// 统计数据
	var totalCredits float32
	var totalGradePoints float32
	courseCount := len(grades)
	passedCourses := 0
	var scores []float32

	for _, grade := range grades {
		scores = append(scores, grade.Score)
		if grade.Score >= 60 {
			passedCourses++
		}
		totalCredits += grade.Course.Credits
		totalGradePoints += grade.GradePoint * grade.Course.Credits
	}

	// 计算GPA
	gpa := totalGradePoints / totalCredits

	// 返回统计结果
	return map[string]interface{}{
		"student_id":         studentID,
		"student_name":       grades[0].Student.Username,
		"course_count":       courseCount,
		"passed_courses":     passedCourses,
		"pass_rate":          float64(passedCourses) / float64(courseCount) * 100,
		"total_credits":      totalCredits,
		"gpa":                math.Round(float64(gpa)*100) / 100,
		"grade_distribution": s.calculateDistribution(scores),
	}, nil
}
