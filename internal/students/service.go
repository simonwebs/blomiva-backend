package students

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	repo Repository
}

func NewService(repo Repository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

func CalculateAge(dob time.Time) int {
	now := time.Now()

	age := now.Year() - dob.Year()

	if now.YearDay() < dob.YearDay() {
		age--
	}

	return age
}

func GetAgeGroup(age int) string {
	if age < AgeMinorLimit {
		return AgeGroupMinor
	}

	if age < AgeAdultLimit {
		return AgeGroupTeen
	}

	return AgeGroupAdult
}

func IsMinor(age int) bool {
	return age < AgeMinorLimit
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *StudentService) CreateStudent(ctx context.Context, req CreateStudentRequest) (*Student, error) {
	if req.DateOfBirth == nil {
		return nil, ErrInvalidStudentAge
	}

	schoolID, err := primitive.ObjectIDFromHex(req.SchoolID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	age := CalculateAge(*req.DateOfBirth)
	ageGroup := GetAgeGroup(age)

	studentEmail := NormalizeEmail(req.StudentEmail)
	parentEmail := NormalizeEmail(req.ParentEmail)

	status := StatusPendingApproval
	emailOwner := "student"

	if ageGroup == AgeGroupMinor {
		if parentEmail == "" {
			return nil, ErrParentEmailRequired
		}

		status = StatusPendingVerification
		emailOwner = "parent"
	}

	if ageGroup == AgeGroupTeen || ageGroup == AgeGroupAdult {
		if studentEmail == "" && parentEmail == "" {
			return nil, ErrStudentEmailRequired
		}

		if studentEmail == "" {
			emailOwner = "parent"
		}
	}

	student := &Student{
		ID:       primitive.NewObjectID(),
		SchoolID: schoolID,

		Name:     strings.TrimSpace(req.Name),
		FullName: strings.TrimSpace(req.Name),

		Gender:      strings.TrimSpace(req.Gender),
		DateOfBirth: req.DateOfBirth,
		Age:         age,
		IsMinor:     IsMinor(age),

		Email:        studentEmail,
		StudentEmail: studentEmail,
		PrimaryEmail: studentEmail,

		Guardian: GuardianInfo{
			ParentEmail:        parentEmail,
			MinorGuardianEmail: parentEmail,
			ConsentRequired:    IsMinor(age),
			ConsentGiven:       false,
			GuardianVerified:   false,
		},

		Academic: AcademicInfo{
			Level:        strings.TrimSpace(req.Level),
			Grade:        strings.TrimSpace(req.Grade),
			LearningMode: "physical",
		},

		Verification: VerificationInfo{
			ContactVerified:       false,
			StudentVerified:       false,
			SchoolVerified:        false,
			VerificationSource:    emailOwner,
			VerificationTokenUsed: false,
		},

		Finance: FinanceSnapshot{
			Currency: "GHS",
		},

		Learning: LearningSnapshot{},
		Location: LocationInfo{},

		Status:     status,
		IsApproved: false,
		IsBlocked:  false,
		IsArchived: false,
		Deleted:    false,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return nil, err
	}

	return student, nil
}
