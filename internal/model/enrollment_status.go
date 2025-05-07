package model

type EnrollmentStatus string

const (
	EnrollmentStatusNone     EnrollmentStatus = "none"
	EnrollmentStatusEnrolled EnrollmentStatus = "enrolled"
	EnrollmentStatusDropped  EnrollmentStatus = "dropped"
)
