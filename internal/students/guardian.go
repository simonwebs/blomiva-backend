package students

import "time"

type GuardianInfo struct {
	ParentName         string     `json:"parentName,omitempty" bson:"parentName,omitempty"`
	ParentEmail        string     `json:"parentEmail,omitempty" bson:"parentEmail,omitempty"`
	ParentPhone        string     `json:"parentPhone,omitempty" bson:"parentPhone,omitempty"`
	MinorGuardianEmail string     `json:"minorGuardianEmail,omitempty" bson:"minorGuardianEmail,omitempty"`
	Relationship       string     `json:"relationship,omitempty" bson:"relationship,omitempty"`
	GuardianVerified   bool       `json:"guardianVerified" bson:"guardianVerified"`
	GuardianVerifiedAt *time.Time `json:"guardianVerifiedAt,omitempty" bson:"guardianVerifiedAt,omitempty"`
	ConsentRequired    bool       `json:"consentRequired" bson:"consentRequired"`
	ConsentGiven       bool       `json:"consentGiven" bson:"consentGiven"`
	ConsentGivenAt     *time.Time `json:"consentGivenAt,omitempty" bson:"consentGivenAt,omitempty"`
}
