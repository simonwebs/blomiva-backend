package tenant

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailConfig struct {
	PostmarkToken string
	FromEmail     string
	AppURL        string
	Secret        string
}

type Service interface {
	Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error)
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	List(ctx context.Context, page, limit int64) ([]Tenant, int64, error)
	Update(ctx context.Context, id string, req UpdateTenantRequest) (*Tenant, error)
	Delete(ctx context.Context, id string) error
	VerifySchool(ctx context.Context, tenantID, email, token string) (*Tenant, error)
}

type TenantService struct {
	repo  Repository
	email EmailConfig
}

// Delete implements [Service].
func (s *TenantService) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByID implements [Service].
func (s *TenantService) GetByID(ctx context.Context, id string) (*Tenant, error) {
	panic("unimplemented")
}

// GetBySlug implements [Service].
func (s *TenantService) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	panic("unimplemented")
}

// List implements [Service].
func (s *TenantService) List(ctx context.Context, page int64, limit int64) ([]Tenant, int64, error) {
	panic("unimplemented")
}

// Update implements [Service].
func (s *TenantService) Update(ctx context.Context, id string, req UpdateTenantRequest) (*Tenant, error) {
	panic("unimplemented")
}

func NewService(r Repository, email EmailConfig) *TenantService {
	return &TenantService{
		repo:  r,
		email: email,
	}
}

// ---------------- CREATE ----------------

func (s *TenantService) Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name required")
	}

	now := time.Now().UTC()

	t := &Tenant{
		ID:       primitive.NewObjectID(),
		TenantID: uuid.NewString(),
		Name:     strings.TrimSpace(req.Name),
		Slug:     normalize(req.Slug),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Status:   StatusPending,
		Verification: TenantVerification{
			EmailVerified: false,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	// SAFE async email (FIXED)
	go s.safeEmail(t)

	return t, nil
}

// ---------------- SAFE EMAIL WRAPPER ----------------

func (s *TenantService) safeEmail(t *Tenant) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("email panic recovered:", r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = s.sendVerification(ctx, t)
}
