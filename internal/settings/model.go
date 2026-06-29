package settings

import "time"

type NotificationSettings struct {
	Email               bool `bson:"email" json:"email"`
	Push                bool `bson:"push" json:"push"`
	SMS                 bool `bson:"sms" json:"sms"`
	SchoolUpdates       bool `bson:"schoolUpdates" json:"schoolUpdates"`
	BillingAlerts       bool `bson:"billingAlerts" json:"billingAlerts"`
	SecurityAlerts      bool `bson:"securityAlerts" json:"securityAlerts"`
	MarketingEmails     bool `bson:"marketingEmails" json:"marketingEmails"`
	ParentUpdates       bool `bson:"parentUpdates" json:"parentUpdates"`
	AttendanceAlerts    bool `bson:"attendanceAlerts" json:"attendanceAlerts"`
	AssignmentAlerts    bool `bson:"assignmentAlerts" json:"assignmentAlerts"`
	FeePaymentReminders bool `bson:"feePaymentReminders" json:"feePaymentReminders"`
}

type PrivacySettings struct {
	ProfileVisibility string `bson:"profileVisibility" json:"profileVisibility"`
	ShowEmail         bool   `bson:"showEmail" json:"showEmail"`
	ShowPhone         bool   `bson:"showPhone" json:"showPhone"`
	ShowOnlineStatus  bool   `bson:"showOnlineStatus" json:"showOnlineStatus"`
	AllowMessages     bool   `bson:"allowMessages" json:"allowMessages"`
}

type SecuritySettings struct {
	TwoFactorEnabled bool      `bson:"twoFactorEnabled" json:"twoFactorEnabled"`
	LoginAlerts      bool      `bson:"loginAlerts" json:"loginAlerts"`
	LastPasswordAt   time.Time `bson:"lastPasswordAt,omitempty" json:"lastPasswordAt,omitempty"`
}

type ProductSettings struct {
	DefaultSchoolID string `bson:"defaultSchoolId" json:"defaultSchoolId"`
	DefaultRole     string `bson:"defaultRole" json:"defaultRole"`
	DashboardMode   string `bson:"dashboardMode" json:"dashboardMode"`
	CompactUI       bool   `bson:"compactUi" json:"compactUi"`
	ReduceMotion    bool   `bson:"reduceMotion" json:"reduceMotion"`
}

type AccountDeletion struct {
	IsDeleted          bool       `bson:"isDeleted" json:"isDeleted"`
	DeleteRequestedAt  *time.Time `bson:"deleteRequestedAt,omitempty" json:"deleteRequestedAt,omitempty"`
	DeleteScheduledAt  *time.Time `bson:"deleteScheduledAt,omitempty" json:"deleteScheduledAt,omitempty"`
	RestoredAt         *time.Time `bson:"restoredAt,omitempty" json:"restoredAt,omitempty"`
	DeletionReason     string     `bson:"deletionReason,omitempty" json:"deletionReason,omitempty"`
	DeletionGraceDays  int        `bson:"deletionGraceDays" json:"deletionGraceDays"`
	DeletionCancelNote string     `bson:"deletionCancelNote,omitempty" json:"deletionCancelNote,omitempty"`
}

type UserSettings struct {
	UserID        string               `bson:"userId" json:"userId"`
	Theme         string               `bson:"theme" json:"theme"`
	Language      string               `bson:"language" json:"language"`
	Timezone      string               `bson:"timezone" json:"timezone"`
	Notifications NotificationSettings `bson:"notifications" json:"notifications"`
	Privacy       PrivacySettings      `bson:"privacy" json:"privacy"`
	Security      SecuritySettings     `bson:"security" json:"security"`
	Product       ProductSettings      `bson:"product" json:"product"`
	Account       AccountDeletion      `bson:"account" json:"account"`
	CreatedAt     time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type UpdateSettingsRequest struct {
	Theme         *string               `json:"theme"`
	Language      *string               `json:"language"`
	Timezone      *string               `json:"timezone"`
	Notifications *NotificationSettings `json:"notifications"`
	Privacy       *PrivacySettings      `json:"privacy"`
	Security      *SecuritySettings     `json:"security"`
	Product       *ProductSettings      `json:"product"`
}

type DeleteAccountRequest struct {
	Reason string `json:"reason"`
}

type RestoreAccountRequest struct {
	Note string `json:"note"`
}
