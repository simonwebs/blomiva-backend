package contact

import "time"

type MessageStatus string

const (
	StatusNew      MessageStatus = "new"
	StatusRead     MessageStatus = "read"
	StatusResolved MessageStatus = "resolved"
)

type ContactMessage struct {
	ID          string        `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Email       string        `bson:"email" json:"email"`
	Description string        `bson:"description" json:"description"`
	Status      MessageStatus `bson:"status" json:"status"`
	IP          string        `bson:"ip,omitempty" json:"ip,omitempty"`
	UserAgent   string        `bson:"userAgent,omitempty" json:"userAgent,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type CreateContactRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Description string `json:"description" binding:"required"`
}

type UpdateContactStatusRequest struct {
	Status MessageStatus `json:"status" binding:"required"`
}
