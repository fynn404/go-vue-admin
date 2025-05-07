-- 创建用户表
CREATE TABLE IF NOT EXISTS users_tab (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CHECK (role IN ('student', 'teacher', 'admin'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建课程表
CREATE TABLE IF NOT EXISTS courses_tab (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    teacher_id BIGINT UNSIGNED NOT NULL,
    credits INT NOT NULL,
    current_enrolled INT NOT NULL DEFAULT 0,
    capacity INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CHECK (status IN ('open', 'closed', 'archived'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户表索引
CREATE UNIQUE INDEX idx_users_username ON users_tab(username);
CREATE INDEX idx_users_role ON users_tab(role);
CREATE INDEX idx_users_deleted_at ON users_tab(deleted_at);

-- 创建课程表索引
CREATE INDEX idx_courses_teacher_id ON courses_tab(teacher_id);
CREATE INDEX idx_courses_status ON courses_tab(status);
CREATE INDEX idx_courses_deleted_at ON courses_tab(deleted_at);

-- 创建成绩表
CREATE TABLE IF NOT EXISTS grade_tab (
                                         id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                         course_id BIGINT UNSIGNED NOT NULL,
                                         student_id BIGINT UNSIGNED NOT NULL,
                                         teacher_id BIGINT UNSIGNED NOT NULL,
                                         score DECIMAL(5,2) NOT NULL,
    grade_point DECIMAL(3,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    comment TEXT,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CHECK (status IN ('draft', 'published'))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建成绩历史记录表
CREATE TABLE IF NOT EXISTS grade_history_tab (
                                                 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                                 grade_id BIGINT UNSIGNED NOT NULL,
                                                 course_id BIGINT UNSIGNED NOT NULL,
                                                 student_id BIGINT UNSIGNED NOT NULL,
                                                 teacher_id BIGINT UNSIGNED NOT NULL,
                                                 change_type VARCHAR(20) NOT NULL,
    old_score DECIMAL(5,2),
    new_score DECIMAL(5,2),
    old_comment TEXT,
    new_comment TEXT,
    reason TEXT,
    operated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CHECK (change_type IN ('create', 'update', 'delete'))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建选课记录表
CREATE TABLE IF NOT EXISTS enrollment_tab (
                                              id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                              student_id BIGINT UNSIGNED NOT NULL,
                                              course_id BIGINT UNSIGNED NOT NULL,
                                              grade DECIMAL(5,2),
    status VARCHAR(20) NOT NULL DEFAULT 'enrolled',
    enrolled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dropped_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CHECK (status IN ('none', 'enrolled', 'dropped'))
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建索引
-- 成绩表索引
CREATE INDEX idx_grades_course_id ON grade_tab(course_id);
CREATE INDEX idx_grades_student_id ON grade_tab(student_id);
CREATE INDEX idx_grades_teacher_id ON grade_tab(teacher_id);
CREATE INDEX idx_grades_status ON grade_tab(status);
CREATE INDEX idx_grades_deleted_at ON grade_tab(deleted_at);

-- 成绩历史记录表索引
CREATE INDEX idx_grade_histories_grade_id ON grade_history_tab(grade_id);
CREATE INDEX idx_grade_histories_course_id ON grade_history_tab(course_id);
CREATE INDEX idx_grade_histories_student_id ON grade_history_tab(student_id);
CREATE INDEX idx_grade_histories_teacher_id ON grade_history_tab(teacher_id);
CREATE INDEX idx_grade_histories_change_type ON grade_history_tab(change_type);
CREATE INDEX idx_grade_histories_operated_at ON grade_history_tab(operated_at);
CREATE INDEX idx_grade_histories_deleted_at ON grade_history_tab(deleted_at);

-- 选课记录表索引
CREATE INDEX idx_enrollments_student_id ON enrollment_tab(student_id);
CREATE INDEX idx_enrollments_course_id ON enrollment_tab(course_id);
CREATE INDEX idx_enrollments_status ON enrollment_tab(status);
CREATE INDEX idx_enrollments_deleted_at ON enrollment_tab(deleted_at);

-- 创建唯一约束
CREATE UNIQUE INDEX idx_enrollments_student_course ON enrollment_tab(student_id, course_id, status);
CREATE UNIQUE INDEX idx_grades_student_course ON grade_tab(student_id, course_id);
