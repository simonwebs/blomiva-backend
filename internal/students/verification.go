package students

import "time"

type VerificationInfo struct {
	ContactVerified       bool       `json:"contactVerified" bson:"contactVerified"`
	ContactVerifiedAt     *time.Time `json:"contactVerifiedAt,omitempty" bson:"contactVerifiedAt,omitempty"`
	StudentVerified       bool       `json:"studentVerified" bson:"studentVerified"`
	StudentVerifiedAt     *time.Time `json:"studentVerifiedAt,omitempty" bson:"studentVerifiedAt,omitempty"`
	SchoolVerified        bool       `json:"schoolVerified" bson:"schoolVerified"`
	SchoolVerifiedAt      *time.Time `json:"schoolVerifiedAt,omitempty" bson:"schoolVerifiedAt,omitempty"`
	VerificationSource    string     `json:"verificationSource,omitempty" bson:"verificationSource,omitempty"`
	VerificationTokenHash string     `json:"-" bson:"verificationTokenHash,omitempty"`
	VerificationExpiresAt *time.Time `json:"verificationExpiresAt,omitempty" bson:"verificationExpiresAt,omitempty"`
	VerificationTokenUsed bool       `json:"verificationTokenUsed" bson:"verificationTokenUsed"`
}
