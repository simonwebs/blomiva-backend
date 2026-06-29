package contact

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func clean(value string) string {
	return strings.TrimSpace(value)
}

func (s *Service) Create(
	ctx context.Context,
	req CreateContactRequest,
	ip string,
	userAgent string,
) (*ContactMessage, error) {
	name := clean(req.Name)
	email := strings.ToLower(clean(req.Email))
	description := clean(req.Description)

	if len(name) < 2 {
		return nil, errors.New("name is too short")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email")
	}

	if len(description) < 10 {
		return nil, errors.New("description is too short")
	}

	if len(description) > 5000 {
		return nil, errors.New("description is too long")
	}

	now := time.Now().UTC()

	message := &ContactMessage{
		ID:          uuid.NewString(),
		Name:        name,
		Email:       email,
		Description: description,
		Status:      StatusNew,
		IP:          ip,
		UserAgent:   userAgent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Service) List(ctx context.Context, limit int64) ([]ContactMessage, error) {
	return s.repo.List(ctx, limit)
}

func (s *Service) Get(ctx context.Context, id string) (*ContactMessage, error) {
	if clean(id) == "" {
		return nil, errors.New("message id is required")
	}

	message, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.New("message not found")
	}

	return message, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status MessageStatus) error {
	if status != StatusNew && status != StatusRead && status != StatusResolved {
		return errors.New("invalid status")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}
