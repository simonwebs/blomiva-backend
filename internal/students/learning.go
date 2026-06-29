package students

import "time"

type LearningSnapshot struct {
	AttendanceRate   float64    `json:"attendanceRate" bson:"attendanceRate"`
	ProgressRate     float64    `json:"progressRate" bson:"progressRate"`
	AverageScore     float64    `json:"averageScore" bson:"averageScore"`
	LessonsCompleted int64      `json:"lessonsCompleted" bson:"lessonsCompleted"`
	AssignmentsDone  int64      `json:"assignmentsDone" bson:"assignmentsDone"`
	ExamsTaken       int64      `json:"examsTaken" bson:"examsTaken"`
	LastSeenAt       *time.Time `json:"lastSeenAt,omitempty" bson:"lastSeenAt,omitempty"`
}
