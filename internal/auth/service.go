package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo       *Repository
	jwtService *JWTService
}

func NewService(repo *Repository, jwtService *JWTService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func normalizeAuthRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	role = strings.ReplaceAll(role, "_", "-")

	switch role {
	case "superadmin", "super-admin":
		return "super-admin"
	case "schooladmin", "school-admin":
		return "school-admin"
	case "":
		return "user"
	default:
		return role
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	email := normalizeEmail(req.Email)
	role := normalizeAuthRole(req.Role)

	if email == "" {
		return nil, errors.New("email required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("password required")
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, email, hash, role, req.SchoolID)
	if err != nil {
		return nil, err
	}

	role = normalizeAuthRole(user.Role)

	token, err := s.jwtService.Generate(user.ID, email, role, user.SchoolID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: AuthUserResponse{
			ID:       user.ID,
			Email:    email,
			Role:     role,
			SchoolID: user.SchoolID,
		},
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	email := normalizeEmail(req.Email)

	if email == "" {
		return nil, errors.New("email required")
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if strings.TrimSpace(user.PasswordHash) == "" {
		return nil, errors.New("password login is not enabled for this account")
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	role := normalizeAuthRole(user.Role)

	token, err := s.jwtService.Generate(user.ID, email, role, user.SchoolID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: AuthUserResponse{
			ID:       user.ID,
			Email:    email,
			Role:     role,
			SchoolID: user.SchoolID,
		},
	}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*AuthUserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	email := ""
	if len(user.Emails) > 0 {
		email = normalizeEmail(user.Emails[0].Address)
	}

	return &AuthUserResponse{
		ID:       user.ID,
		Email:    email,
		Role:     normalizeAuthRole(user.Role),
		SchoolID: user.SchoolID,
	}, nil
}

func (s *Service) SeedSuperAdmin(ctx context.Context) error {
	email := normalizeEmail(os.Getenv("SUPER_ADMIN_EMAIL"))
	password := strings.TrimSpace(os.Getenv("SUPER_ADMIN_PASSWORD"))

	if email == "" {
		return errors.New("SUPER_ADMIN_EMAIL missing")
	}

	if password == "" {
		return errors.New("SUPER_ADMIN_PASSWORD missing")
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	if existing != nil {
		return nil
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.repo.CreateUser(ctx, email, hash, "super-admin", "")
	return err
}
