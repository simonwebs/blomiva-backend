package students

import "errors"

var (
	ErrInvalidStudentAge    = errors.New("invalid student age")
	ErrParentEmailRequired  = errors.New("parent email is required for students under 13")
	ErrStudentEmailRequired = errors.New("student email or parent email is required")
)
