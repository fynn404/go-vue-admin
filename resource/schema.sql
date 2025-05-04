-- 创建课程管理系统数据库
-- 使用UTF8MB4字符集以支持完整的Unicode字符（包括表情符号）
CREATE DATABASE IF NOT EXISTS course_admin 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE course_admin;

-- 用户表：存储所有类型用户信息（管理员、教师、学生）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,           -- 用户ID，自增主键
    username VARCHAR(50) NOT NULL UNIQUE,                    -- 用户名，唯一约束
    password VARCHAR(255) NOT NULL,                         -- 密码，使用bcrypt加密存储
    role ENUM('admin', 'teacher', 'student') NOT NULL,      -- 用户角色：管理员、教师、学生
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
    deleted_at TIMESTAMP NULL,                             -- 软删除时间戳
    INDEX idx_username (username),                         -- 用户名索引，用于登录查询
    INDEX idx_role (role),                                -- 角色索引，用于角色过滤
    INDEX idx_deleted_at (deleted_at)                      -- 软删除索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 课程表：存储课程信息
CREATE TABLE IF NOT EXISTS courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,           -- 课程ID，自增主键
    name VARCHAR(100) NOT NULL,                            -- 课程名称
    description TEXT,                                      -- 课程描述
    teacher_id BIGINT UNSIGNED NOT NULL,                   -- 教师ID，关联users表
    credits DECIMAL(3,1) NOT NULL,                         -- 学分，支持0.5学分
    capacity INT NOT NULL,                                 -- 课程容量
    current_enrolled INT NOT NULL DEFAULT 0,               -- 当前已选人数
    status ENUM('open', 'closed') NOT NULL DEFAULT 'open', -- 课程状态：开放/关闭
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
    deleted_at TIMESTAMP NULL,                             -- 软删除时间戳
    FOREIGN KEY (teacher_id) REFERENCES users(id),         -- 教师外键约束
    INDEX idx_teacher_id (teacher_id),                     -- 教师ID索引，用于查询教师的课程
    INDEX idx_status (status),                            -- 状态索引，用于筛选开放课程
    INDEX idx_deleted_at (deleted_at)                      -- 软删除索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 选课记录表：存储学生选课信息
CREATE TABLE IF NOT EXISTS enrollments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,           -- 选课记录ID，自增主键
    student_id BIGINT UNSIGNED NOT NULL,                   -- 学生ID，关联users表
    course_id BIGINT UNSIGNED NOT NULL,                    -- 课程ID，关联courses表
    status ENUM('active', 'dropped', 'complete') NOT NULL DEFAULT 'active',  -- 选课状态：进行中/已退课/已完成
    grade DECIMAL(5,2) NULL,                               -- 成绩，百分制（0-100）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
    deleted_at TIMESTAMP NULL,                             -- 软删除时间戳
    FOREIGN KEY (student_id) REFERENCES users(id),         -- 学生外键约束
    FOREIGN KEY (course_id) REFERENCES courses(id),        -- 课程外键约束
    UNIQUE KEY uk_student_course (student_id, course_id),  -- 学生和课程的唯一组合，防止重复选课
    INDEX idx_student_id (student_id),                     -- 学生ID索引，用于查询学生的选课记录
    INDEX idx_course_id (course_id),                       -- 课程ID索引，用于查询课程的选课记录
    INDEX idx_status (status),                            -- 状态索引，用于筛选选课状态
    INDEX idx_deleted_at (deleted_at)                      -- 软删除索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 成绩历史记录表：记录成绩变更历史
CREATE TABLE IF NOT EXISTS grade_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,           -- 历史记录ID，自增主键
    enrollment_id BIGINT UNSIGNED NOT NULL,                -- 选课记录ID，关联enrollments表
    teacher_id BIGINT UNSIGNED NOT NULL,                   -- 教师ID，关联users表
    grade DECIMAL(5,2) NOT NULL,                           -- 成绩，百分制（0-100）
    comment TEXT,                                          -- 成绩变更说明
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 变更时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
    deleted_at TIMESTAMP NULL,                             -- 软删除时间戳
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id), -- 选课记录外键约束
    FOREIGN KEY (teacher_id) REFERENCES users(id),         -- 教师外键约束
    INDEX idx_enrollment_id (enrollment_id),               -- 选课记录ID索引，用于查询成绩历史
    INDEX idx_teacher_id (teacher_id),                     -- 教师ID索引，用于查询教师的操作记录
    INDEX idx_timestamp (timestamp),                       -- 时间戳索引，用于时间范围查询
    INDEX idx_deleted_at (deleted_at)                      -- 软删除索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建初始管理员用户
-- 用户名：admin
-- 密码：admin123（使用bcrypt加密）
-- 角色：admin
INSERT INTO users (username, password, role) VALUES 
('admin', '$2a$10$GVvZmMv6P7AGYz7mXXKjBuGHHZ.DQF3FZIgkYqQUc2NWIT0GQlWBe', 'admin'); 