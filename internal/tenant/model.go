package tenant

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	TenantID   string `json:"tenantId" bson:"tenantId"`
	Slug       string `json:"slug" bson:"slug"`
	Name       string `json:"name" bson:"name"`
	LegalName  string `json:"legalName,omitempty" bson:"legalName,omitempty"`
	SchoolCode string `json:"schoolCode,omitempty" bson:"schoolCode,omitempty"`
	Domain     string `json:"domain,omitempty" bson:"domain,omitempty"`

	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
	Website     string `json:"website,omitempty" bson:"website,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`

	Logo    MediaAsset  `json:"logo,omitempty" bson:"logo,omitempty"`
	Banner  MediaAsset  `json:"banner,omitempty" bson:"banner,omitempty"`
	Address Address     `json:"address" bson:"address"`
	Geo     GeoLocation `json:"geo" bson:"geo"`

	Owner   OwnerInfo     `json:"owner" bson:"owner"`
	Consent TenantConsent `json:"consent" bson:"consent"`

	Subscription Subscription   `json:"subscription" bson:"subscription"`
	Billing      Billing        `json:"billing" bson:"billing"`
	Features     FeatureFlags   `json:"features" bson:"features"`
	Settings     TenantSettings `json:"settings" bson:"settings"`

	Verification TenantVerification `json:"verification" bson:"verification"`

	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`

	Status   string `json:"status" bson:"status"`
	Active   bool   `json:"active" bson:"active"`
	Deleted  bool   `json:"deleted" bson:"deleted"`
	Archived bool   `json:"archived" bson:"archived"`

	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type TenantVerification struct {
	EmailVerified bool       `json:"emailVerified" bson:"emailVerified"`
	VerifiedAt    *time.Time `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"`
}

type TenantConsent struct {
	Accepted           bool      `json:"accepted" bson:"accepted"`
	VerificationEmails bool      `json:"verificationEmails" bson:"verificationEmails"`
	AccountUpdates     bool      `json:"accountUpdates" bson:"accountUpdates"`
	PrivacyPolicy      bool      `json:"privacyPolicy" bson:"privacyPolicy"`
	Terms              bool      `json:"terms" bson:"terms"`
	AcceptedAt         time.Time `json:"acceptedAt" bson:"acceptedAt"`
	IPAddress          string    `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	UserAgent          string    `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
}

type Address struct {
	Country string `json:"country,omitempty" bson:"country,omitempty"`
	Region  string `json:"region,omitempty" bson:"region,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	Town    string `json:"town,omitempty" bson:"town,omitempty"`
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
}

type GeoLocation struct {
	Latitude    float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	TimeZone    string  `json:"timeZone,omitempty" bson:"timeZone,omitempty"`
	CountryCode string  `json:"countryCode,omitempty" bson:"countryCode,omitempty"`
}

type OwnerInfo struct {
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`
	Phone  string             `json:"phone,omitempty" bson:"phone,omitempty"`
}

type MediaAsset struct {
	URL       string     `json:"url,omitempty" bson:"url,omitempty"`
	Key       string     `json:"key,omitempty" bson:"key,omitempty"`
	MimeType  string     `json:"mimeType,omitempty" bson:"mimeType,omitempty"`
	Size      int64      `json:"size,omitempty" bson:"size,omitempty"`
	Width     int        `json:"width,omitempty" bson:"width,omitempty"`
	Height    int        `json:"height,omitempty" bson:"height,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type Logo = MediaAsset
type Banner = MediaAsset

type Subscription struct {
	Plan         string     `json:"plan" bson:"plan"`
	Status       string     `json:"status" bson:"status"`
	BillingCycle string     `json:"billingCycle" bson:"billingCycle"`
	Currency     string     `json:"currency" bson:"currency"`
	Amount       float64    `json:"amount" bson:"amount"`
	Trial        bool       `json:"trial" bson:"trial"`
	TrialEndsAt  *time.Time `json:"trialEndsAt,omitempty" bson:"trialEndsAt,omitempty"`
	AutoRenew    bool       `json:"autoRenew" bson:"autoRenew"`
	RenewsAt     *time.Time `json:"renewsAt,omitempty" bson:"renewsAt,omitempty"`
}

type Billing struct {
	CustomerID        string `json:"customerId,omitempty" bson:"customerId,omitempty"`
	Provider          string `json:"provider,omitempty" bson:"provider,omitempty"`
	BillingEmail      string `json:"billingEmail,omitempty" bson:"billingEmail,omitempty"`
	TaxID             string `json:"taxId,omitempty" bson:"taxId,omitempty"`
	DefaultCurrency   string `json:"defaultCurrency,omitempty" bson:"defaultCurrency,omitempty"`
	LastPaymentStatus string `json:"lastPaymentStatus,omitempty" bson:"lastPaymentStatus,omitempty"`
}

type FeatureFlags struct {
	Students          bool `json:"students" bson:"students"`
	Teachers          bool `json:"teachers" bson:"teachers"`
	Parents           bool `json:"parents" bson:"parents"`
	Attendance        bool `json:"attendance" bson:"attendance"`
	Exams             bool `json:"exams" bson:"exams"`
	ReportCards       bool `json:"reportCards" bson:"reportCards"`
	Accounting        bool `json:"accounting" bson:"accounting"`
	Messaging         bool `json:"messaging" bson:"messaging"`
	ParentPortal      bool `json:"parentPortal" bson:"parentPortal"`
	StudentPortal     bool `json:"studentPortal" bson:"studentPortal"`
	TeacherPortal     bool `json:"teacherPortal" bson:"teacherPortal"`
	Email             bool `json:"email" bson:"email"`
	PushNotifications bool `json:"pushNotifications" bson:"pushNotifications"`
}

type TenantSettings struct {
	DefaultLanguage   string `json:"defaultLanguage" bson:"defaultLanguage"`
	DefaultCurrency   string `json:"defaultCurrency" bson:"defaultCurrency"`
	TimeZone          string `json:"timeZone" bson:"timeZone"`
	DateFormat        string `json:"dateFormat" bson:"dateFormat"`
	TimeFormat        string `json:"timeFormat" bson:"timeFormat"`
	Theme             string `json:"theme" bson:"theme"`
	AllowRegistration bool   `json:"allowRegistration" bson:"allowRegistration"`
	RequireApproval   bool   `json:"requireApproval" bson:"requireApproval"`
}
