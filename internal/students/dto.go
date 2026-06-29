package students

import "time"

type CreateStudentRequest struct {
	SchoolID string `json:"schoolId"`

	Name        string     `json:"name"`
	Gender      string     `json:"gender"`
	DateOfBirth *time.Time `json:"dateOfBirth"`

	StudentEmail string `json:"studentEmail"`
	ParentEmail  string `json:"parentEmail"`

	Level   string `json:"level"`
	Grade   string `json:"grade"`
	ClassID string `json:"classId"`
}
