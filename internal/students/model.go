package students

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SchoolID primitive.ObjectID `json:"schoolId" bson:"schoolId"`

	UserID             *primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	StudentUserID      *primitive.ObjectID `json:"studentUserId,omitempty" bson:"studentUserId,omitempty"`
	ParentUserID       *primitive.ObjectID `json:"parentUserId,omitempty" bson:"parentUserId,omitempty"`
	SchoolMembershipID *primitive.ObjectID `json:"schoolMembershipId,omitempty" bson:"schoolMembershipId,omitempty"`

	StudentNumber string `json:"studentNumber" bson:"studentNumber"`
	Name          string `json:"name" bson:"name"`
	FullName      string `json:"fullName" bson:"fullName"`

	Gender      string     `json:"gender" bson:"gender"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty"`
	Age         int        `json:"age" bson:"age"`
	IsMinor     bool       `json:"isMinor" bson:"isMinor"`

	Email        string `json:"email,omitempty" bson:"email,omitempty"`
	StudentEmail string `json:"studentEmail,omitempty" bson:"studentEmail,omitempty"`
	PrimaryEmail string `json:"primaryEmail,omitempty" bson:"primaryEmail,omitempty"`

	PhoneNumber string `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Address     string `json:"address,omitempty" bson:"address,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty" bson:"photoUrl,omitempty"`

	Academic     AcademicInfo     `json:"academic" bson:"academic"`
	Guardian     GuardianInfo     `json:"guardian" bson:"guardian"`
	Verification VerificationInfo `json:"verification" bson:"verification"`
	Finance      FinanceSnapshot  `json:"finance" bson:"finance"`
	Learning     LearningSnapshot `json:"learning" bson:"learning"`
	Location     LocationInfo     `json:"location" bson:"location"`

	Status string `json:"status" bson:"status"`

	IsApproved bool `json:"isApproved" bson:"isApproved"`
	IsBlocked  bool `json:"isBlocked" bson:"isBlocked"`
	IsArchived bool `json:"isArchived" bson:"isArchived"`
	Deleted    bool `json:"deleted" bson:"deleted"`

	CreatedBy  *primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedBy  *primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	ApprovedBy *primitive.ObjectID `json:"approvedBy,omitempty" bson:"approvedBy,omitempty"`

	ApprovedAt *time.Time `json:"approvedAt,omitempty" bson:"approvedAt,omitempty"`
	RejectedAt *time.Time `json:"rejectedAt,omitempty" bson:"rejectedAt,omitempty"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" bson:"archivedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
