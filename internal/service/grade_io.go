package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	models "github.com/fynn404/go-vue-admin/internal/model"
	"github.com/fynn404/go-vue-admin/internal/repository"
)

// ValidationError 数据验证错误
type ValidationError struct {
	Row     int    // 数据行号
	Field   string // 字段名
	Message string // 错误信息
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("第 %d 行 %s: %s", e.Row, e.Field, e.Message)
}

// ValidationErrors 数据验证错误集合
type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	result := fmt.Sprintf("发现 %d 个错误:\n", len(e))
	for _, err := range e {
		result += err.Error() + "\n"
	}
	return result
}

// GradeIO 成绩导入导出服务
type GradeIO struct {
	gradeRepo  repository.GradeRepository
	userRepo   repository.UserRepository
	courseRepo repository.CourseRepository
}

// NewGradeIO 创建成绩导入导出服务实例
func NewGradeIO(gradeRepo repository.GradeRepository, userRepo repository.UserRepository, courseRepo repository.CourseRepository) *GradeIO {
	return &GradeIO{
		gradeRepo:  gradeRepo,
		userRepo:   userRepo,
		courseRepo: courseRepo,
	}
}

// GradeExportData 成绩导出数据结构
type GradeExportData struct {
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	CourseID    string    `json:"course_id"`
	CourseName  string    `json:"course_name"`
	TeacherName string    `json:"teacher_name"`
	Score       float32   `json:"score"`
	GradePoint  float32   `json:"grade_point"`
	Status      string    `json:"status"`
	Comment     string    `json:"comment"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

// GradeImportData 成绩导入数据结构
type GradeImportData struct {
	StudentID string  `json:"student_id"`
	CourseID  string  `json:"course_id"`
	TeacherID string  `json:"teacher_id"`
	Score     float32 `json:"score"`
	Comment   string  `json:"comment"`
}

// ExportToCSV 导出成绩到CSV文件
func (s *GradeIO) ExportToCSV(ctx context.Context, writer io.Writer, courseID uint) error {
	// 获取课程成绩数据
	grades, _, err := s.gradeRepo.ListByCourse(ctx, courseID, 0, 1000)
	if err != nil {
		return fmt.Errorf("获取成绩数据失败: %w", err)
	}

	// 创建CSV writer
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// 写入表头
	headers := []string{"学号", "学生姓名", "课程编号", "课程名称", "教师姓名", "分数", "绩点", "状态", "评语", "发布时间"}
	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}

	// 写入数据
	for _, grade := range grades {
		publishedAt := ""
		if grade.PublishedAt != nil {
			publishedAt = grade.PublishedAt.Format("2006-01-02 15:04:05")
		}

		row := []string{
			fmt.Sprintf("%d", grade.StudentID),
			grade.Student.Username,
			fmt.Sprintf("%d", grade.CourseID),
			grade.Course.Name,
			grade.Teacher.Username,
			fmt.Sprintf("%.1f", grade.Score),
			fmt.Sprintf("%.1f", grade.GradePoint),
			string(grade.Status),
			grade.Comment,
			publishedAt,
		}

		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("写入数据行失败: %w", err)
		}
	}

	return nil
}

// ExportToExcel 导出成绩到Excel文件
func (s *GradeIO) ExportToExcel(ctx context.Context, writer io.Writer, courseID uint) error {
	// 获取课程成绩数据
	grades, _, err := s.gradeRepo.ListByCourse(ctx, courseID, 0, 1000)
	if err != nil {
		return fmt.Errorf("获取成绩数据失败: %w", err)
	}

	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置表头
	headers := []string{"学号", "学生姓名", "课程编号", "课程名称", "教师姓名", "分数", "绩点", "状态", "评语", "发布时间"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// 写入数据
	for i, grade := range grades {
		row := i + 2 // 从第二行开始写入数据
		publishedAt := ""
		if grade.PublishedAt != nil {
			publishedAt = grade.PublishedAt.Format("2006-01-02 15:04:05")
		}

		data := []interface{}{
			grade.StudentID,
			grade.Student.Username,
			grade.CourseID,
			grade.Course.Name,
			grade.Teacher.Username,
			grade.Score,
			grade.GradePoint,
			string(grade.Status),
			grade.Comment,
			publishedAt,
		}

		for j, value := range data {
			cell := fmt.Sprintf("%c%d", 'A'+j, row)
			f.SetCellValue("Sheet1", cell, value)
		}
	}

	// 设置列宽
	for i := 0; i < len(headers); i++ {
		col := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth("Sheet1", col, col, 15)
	}

	// 写入到writer
	return f.Write(writer)
}

// validateGradeData 验证成绩数据
func (s *GradeIO) validateGradeData(ctx context.Context, grades []*models.Grade) ValidationErrors {
	var errors ValidationErrors

	for i, grade := range grades {
		rowNum := i + 2 // 加2是因为第1行是表头

		// 验证学生ID
		if student, err := s.userRepo.FindByID(ctx, grade.StudentID); err != nil || student == nil {
			errors = append(errors, &ValidationError{
				Row:     rowNum,
				Field:   "学号",
				Message: fmt.Sprintf("学号 %d 不存在", grade.StudentID),
			})
		}

		// 验证课程ID
		if course, err := s.courseRepo.FindByID(ctx, grade.CourseID); err != nil || course == nil {
			errors = append(errors, &ValidationError{
				Row:     rowNum,
				Field:   "课程编号",
				Message: fmt.Sprintf("课程编号 %d 不存在", grade.CourseID),
			})
		}

		// 验证教师ID
		if teacher, err := s.userRepo.FindByID(ctx, grade.TeacherID); err != nil || teacher == nil {
			errors = append(errors, &ValidationError{
				Row:     rowNum,
				Field:   "教师编号",
				Message: fmt.Sprintf("教师编号 %d 不存在", grade.TeacherID),
			})
		}

		// 验证分数范围
		if grade.Score < 0 || grade.Score > 100 {
			errors = append(errors, &ValidationError{
				Row:     rowNum,
				Field:   "分数",
				Message: fmt.Sprintf("分数 %.1f 超出有效范围(0-100)", grade.Score),
			})
		}
	}

	return errors
}

// ImportFromCSV 从CSV文件导入成绩
func (s *GradeIO) ImportFromCSV(ctx context.Context, reader io.Reader) error {
	// 创建CSV reader
	csvReader := csv.NewReader(reader)

	// 跳过表头
	if _, err := csvReader.Read(); err != nil {
		return fmt.Errorf("读取CSV表头失败: %w", err)
	}

	var grades []*models.Grade

	// 读取数据行
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取CSV数据行失败: %w", err)
		}

		// 解析数据
		studentID, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			return fmt.Errorf("无效的学号格式: %s", record[0])
		}

		courseID, err := strconv.ParseUint(record[2], 10, 32)
		if err != nil {
			return fmt.Errorf("无效的课程编号格式: %s", record[2])
		}

		teacherID, err := strconv.ParseUint(record[4], 10, 32)
		if err != nil {
			return fmt.Errorf("无效的教师编号格式: %s", record[4])
		}

		score, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return fmt.Errorf("无效的分数格式: %s", record[5])
		}

		grade := &models.Grade{
			StudentID: uint(studentID),
			CourseID:  uint(courseID),
			TeacherID: uint(teacherID),
			Score:     float32(score),
			Comment:   record[8],
			Status:    models.GradeStatusDraft,
		}

		grades = append(grades, grade)
	}

	// 验证数据
	if errors := s.validateGradeData(ctx, grades); len(errors) > 0 {
		return errors
	}

	// 批量创建成绩记录
	if err := s.gradeRepo.BatchCreate(ctx, grades); err != nil {
		return fmt.Errorf("批量创建成绩记录失败: %w", err)
	}

	return nil
}

// ImportFromExcel 从Excel文件导入成绩
func (s *GradeIO) ImportFromExcel(ctx context.Context, reader io.Reader) error {
	// 读取Excel文件
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	// 获取所有行
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("读取Excel数据失败: %w", err)
	}

	var grades []*models.Grade

	// 跳过表头，处理数据行
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 9 {
			continue
		}

		// 解析数据
		studentID, err := strconv.ParseUint(row[0], 10, 32)
		if err != nil {
			return fmt.Errorf("第 %d 行：无效的学号格式: %s", i+1, row[0])
		}

		courseID, err := strconv.ParseUint(row[2], 10, 32)
		if err != nil {
			return fmt.Errorf("第 %d 行：无效的课程编号格式: %s", i+1, row[2])
		}

		teacherID, err := strconv.ParseUint(row[4], 10, 32)
		if err != nil {
			return fmt.Errorf("第 %d 行：无效的教师编号格式: %s", i+1, row[4])
		}

		score, err := strconv.ParseFloat(row[5], 32)
		if err != nil {
			return fmt.Errorf("第 %d 行：无效的分数格式: %s", i+1, row[5])
		}

		grade := &models.Grade{
			StudentID: uint(studentID),
			CourseID:  uint(courseID),
			TeacherID: uint(teacherID),
			Score:     float32(score),
			Comment:   row[8],
			Status:    models.GradeStatusDraft,
		}

		grades = append(grades, grade)
	}

	// 验证数据
	if errors := s.validateGradeData(ctx, grades); len(errors) > 0 {
		return errors
	}

	// 批量创建成绩记录
	if err := s.gradeRepo.BatchCreate(ctx, grades); err != nil {
		return fmt.Errorf("批量创建成绩记录失败: %w", err)
	}

	return nil
}
