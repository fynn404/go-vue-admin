package model

import "errors"

var (
	// Course related errors
	ErrInvalidCapacity   = errors.New("invalid course capacity")
	ErrCourseUnavailable = errors.New("course is not available for enrollment")
	ErrNoEnrollment      = errors.New("no enrollment to cancel")
	ErrStudentNotFound   = errors.New("student not found in course")
	ErrCourseFull        = errors.New("course has reached maximum capacity")

	// User related errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUsername  = errors.New("username already exists")

	// Enrollment related errors
	ErrAlreadyEnrolled         = errors.New("student already enrolled in course")
	ErrNotEnrolled             = errors.New("student not enrolled in course")
	ErrDuplicateEnrollment     = errors.New("student already enrolled in this course")
	ErrInvalidEnrollmentStatus = errors.New("invalid enrollment status")
	ErrInvalidStudent          = errors.New("invalid student")
	ErrInvalidGrade            = errors.New("invalid grade")
)
