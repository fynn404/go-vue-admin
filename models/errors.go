package models

import "errors"

var (
	// Course related errors
	ErrInvalidCapacity   = errors.New("course capacity cannot be less than current enrolled")
	ErrCourseUnavailable = errors.New("course is not available for enrollment")
	ErrNoEnrollment      = errors.New("no enrollment to remove")
	ErrCourseFull        = errors.New("course has reached maximum capacity")

	// User related errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUsername  = errors.New("username already exists")

	// Enrollment related errors
	ErrAlreadyEnrolled = errors.New("student already enrolled in course")
	ErrNotEnrolled     = errors.New("student not enrolled in course")
)
