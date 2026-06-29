package students

import "go.mongodb.org/mongo-driver/bson/primitive"

type AcademicInfo struct {
	AcademicYear string              `json:"academicYear,omitempty" bson:"academicYear,omitempty"`
	Level        string              `json:"level,omitempty" bson:"level,omitempty"`
	Grade        string              `json:"grade,omitempty" bson:"grade,omitempty"`
	GradeID      *primitive.ObjectID `json:"gradeId,omitempty" bson:"gradeId,omitempty"`
	ClassID      *primitive.ObjectID `json:"classId,omitempty" bson:"classId,omitempty"`
	ClassName    string              `json:"className,omitempty" bson:"className,omitempty"`
	ClassLevel   string              `json:"classLevel,omitempty" bson:"classLevel,omitempty"`
	LearningMode string              `json:"learningMode,omitempty" bson:"learningMode,omitempty"`
}
