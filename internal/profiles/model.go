package profiles

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	Provider    string     `bson:"provider,omitempty" json:"provider,omitempty"`
	Bucket      string     `bson:"bucket,omitempty" json:"bucket,omitempty"`
	Key         string     `bson:"key,omitempty" json:"key,omitempty"`
	StorageKey  string     `bson:"storageKey,omitempty" json:"storageKey,omitempty"`
	PublicID    string     `bson:"publicId,omitempty" json:"publicId,omitempty"`
	URL         string     `bson:"url,omitempty" json:"url,omitempty"`
	PublicURL   string     `bson:"publicUrl,omitempty" json:"publicUrl,omitempty"`
	ContentType string     `bson:"contentType,omitempty" json:"contentType,omitempty"`
	Width       int        `bson:"width,omitempty" json:"width,omitempty"`
	Height      int        `bson:"height,omitempty" json:"height,omitempty"`
	Format      string     `bson:"format,omitempty" json:"format,omitempty"`
	Bytes       int64      `bson:"bytes,omitempty" json:"bytes,omitempty"`
	UploadedAt  *time.Time `bson:"uploadedAt,omitempty" json:"uploadedAt,omitempty"`
}

type Location struct {
	Country string `bson:"country,omitempty" json:"country,omitempty"`
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	Address string `bson:"address,omitempty" json:"address,omitempty"`
}

type SocialLink struct {
	Platform string `bson:"platform" json:"platform"`
	URL      string `bson:"url" json:"url"`
	Label    string `bson:"label,omitempty" json:"label,omitempty"`
}

type Privacy struct {
	ShowEmail       bool `bson:"showEmail" json:"showEmail"`
	ShowPhone       bool `bson:"showPhone" json:"showPhone"`
	ShowLocation    bool `bson:"showLocation" json:"showLocation"`
	ShowSocialLinks bool `bson:"showSocialLinks" json:"showSocialLinks"`
	ShowBio         bool `bson:"showBio" json:"showBio"`
	ShowCustom      bool `bson:"showCustom" json:"showCustom"`
}

type Consent struct {
	Marketing bool       `bson:"marketing" json:"marketing"`
	Email     bool       `bson:"email" json:"email"`
	UpdatedAt *time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type Verification struct {
	IsVerified bool       `bson:"isVerified" json:"isVerified"`
	SentAt     *time.Time `bson:"sentAt,omitempty" json:"sentAt,omitempty"`
	Token      string     `bson:"token,omitempty" json:"-"`
}

type Profile struct {
	ID                   primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	OwnerID              string                 `bson:"ownerId" json:"ownerId"`
	Username             string                 `bson:"username,omitempty" json:"username,omitempty"`
	Slug                 string                 `bson:"slug,omitempty" json:"slug,omitempty"`
	Email                string                 `bson:"email,omitempty" json:"email,omitempty"`
	Phone                string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	Language             string                 `bson:"language,omitempty" json:"language,omitempty"`
	DisplayName          string                 `bson:"displayName,omitempty" json:"displayName,omitempty"`
	Bio                  string                 `bson:"bio,omitempty" json:"bio,omitempty"`
	ProfileImage         *Media                 `bson:"profileImage,omitempty" json:"profileImage,omitempty"`
	Avatar               *Media                 `bson:"avatar,omitempty" json:"avatar,omitempty"`
	AvatarURL            string                 `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty"`
	Banner               *Media                 `bson:"banner,omitempty" json:"banner,omitempty"`
	BannerURL             string                 `bson:"bannerUrl,omitempty" json:"bannerUrl,omitempty"`
	Location             *Location              `bson:"location,omitempty" json:"location,omitempty"`
	SocialLinks          []SocialLink            `bson:"socialLinks,omitempty" json:"socialLinks,omitempty"`
	Privacy              Privacy                `bson:"privacy" json:"privacy"`
	Consent              Consent                `bson:"consent" json:"consent"`
	Verification         Verification           `bson:"verification" json:"verification"`
	Status               string                 `bson:"status" json:"status"`
	IsBlocked            bool                   `bson:"isBlocked" json:"isBlocked"`
	ScheduledForDeletion *time.Time             `bson:"scheduledForDeletion,omitempty" json:"scheduledForDeletion,omitempty"`
	Custom               map[string]interface{} `bson:"custom,omitempty" json:"custom,omitempty"`
	LastActiveAt          *time.Time             `bson:"lastActiveAt,omitempty" json:"lastActiveAt,omitempty"`
	CreatedAt            time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time              `bson:"updatedAt" json:"updatedAt"`
}

type User struct {
	ID                 string                 `bson:"_id" json:"id"`
	Username           string                 `bson:"username,omitempty" json:"username,omitempty"`
	Emails             []UserEmail            `bson:"emails,omitempty" json:"emails,omitempty"`
	Profile            map[string]interface{} `bson:"profile,omitempty" json:"profile,omitempty"`
	MustChangePassword bool                   `bson:"mustChangePassword,omitempty" json:"mustChangePassword,omitempty"`
	CreatedAt          time.Time              `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

type UserEmail struct {
	Address  string `bson:"address" json:"address"`
	Verified bool   `bson:"verified" json:"verified"`
}

type UpdateProfileRequest struct {
	Username    *string                `json:"username,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Phone       *string                `json:"phone,omitempty"`
	Language    *string                `json:"language,omitempty"`
	DisplayName *string                `json:"displayName,omitempty"`
	Bio         *string                `json:"bio,omitempty"`
	Location    *Location              `json:"location,omitempty"`
	SocialLinks *[]SocialLink           `json:"socialLinks,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
}

type UploadImageRequest struct {
	Base64Image string `json:"base64Image" binding:"required"`
}

type SetUserStatusRequest struct {
	OwnerID string `json:"ownerId" binding:"required"`
	Active  bool   `json:"active"`
}