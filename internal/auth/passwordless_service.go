package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordlessPurposeLogin = "login"
	passwordlessCodeLength   = 6
	passwordlessTTLMinutes   = 15
	passwordlessMaxAttempts  = 5
)

type PasswordlessToken struct {
	ID        string     `bson:"_id" json:"id"`
	Email     string     `bson:"email" json:"email"`
	CodeHash  string     `bson:"codeHash" json:"-"`
	Purpose   string     `bson:"purpose" json:"purpose"`
	Attempts  int        `bson:"attempts" json:"attempts"`
	Used      bool       `bson:"used" json:"used"`
	ExpiresAt time.Time  `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
	UsedAt    *time.Time `bson:"usedAt,omitempty" json:"usedAt,omitempty"`
}

type PasswordlessService struct {
	repo       *Repository
	tokens     *mongo.Collection
	emailer    EmailSender
	jwtService *JWTService
	verifyURL  string
	tokenTTL   time.Duration
	maxAttempt int
}

type EmailSender interface {
	Send(to string, subject string, html string) error
}

func NewPasswordlessService(
	db *mongo.Database,
	repo *Repository,
	emailer EmailSender,
	jwtService *JWTService,
) *PasswordlessService {
	verifyURL := strings.TrimSpace(os.Getenv("APP_VERIFY_URL"))
	if verifyURL == "" {
		verifyURL = "http://localhost:8081/auth/passwordless"
	}

	return &PasswordlessService{
		repo:       repo,
		tokens:     db.Collection("passwordless_tokens"),
		emailer:    emailer,
		jwtService: jwtService,
		verifyURL:  verifyURL,
		tokenTTL:   passwordlessTTLMinutes * time.Minute,
		maxAttempt: passwordlessMaxAttempts,
	}
}

func normalizePasswordlessCode(code string) string {
	var b strings.Builder

	for _, r := range strings.TrimSpace(code) {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func validatePasswordlessEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func generatePasswordlessCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func primaryUserEmail(user *User) string {
	if user == nil || len(user.Emails) == 0 {
		return ""
	}

	return normalizeEmail(user.Emails[0].Address)
}

func (s *PasswordlessService) RequestCode(ctx context.Context, email string) error {
	email = normalizeEmail(email)

	if err := validatePasswordlessEmail(email); err != nil {
		return err
	}

	code, err := generatePasswordlessCode()
	if err != nil {
		return errors.New("failed to generate login code")
	}

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to secure login code")
	}

	now := time.Now().UTC()

	_, _ = s.tokens.UpdateMany(
		ctx,
		bson.M{
			"email":   email,
			"purpose": passwordlessPurposeLogin,
			"used":    false,
		},
		bson.M{
			"$set": bson.M{
				"used":      true,
				"updatedAt": now,
			},
		},
	)

	loginToken := PasswordlessToken{
		ID:        uuid.NewString(),
		Email:     email,
		CodeHash:  string(codeHash),
		Purpose:   passwordlessPurposeLogin,
		Attempts:  0,
		Used:      false,
		ExpiresAt: now.Add(s.tokenTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.tokens.InsertOne(ctx, loginToken); err != nil {
		return errors.New("failed to save login code")
	}

	link := fmt.Sprintf("%s?email=%s&code=%s", s.verifyURL, email, code)

	html := fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;line-height:1.6;color:#111827">
			<h2>Blomiva School Login Code</h2>
			<p>Your secure login code is:</p>
			<div style="font-size:32px;font-weight:800;letter-spacing:8px;margin:20px 0">%s</div>
			<p>This code expires in 15 minutes.</p>
			<p>
				<a href="%s" style="display:inline-block;background:#26384F;color:white;padding:12px 18px;border-radius:999px;text-decoration:none;font-weight:700">
					Login to Blomiva School
				</a>
			</p>
			<p style="font-size:12px;color:#6b7280">
				If you did not request this email, you can ignore it.
			</p>
		</div>
	`, code, link)

	return s.emailer.Send(email, "Your Blomiva School login code", html)
}

func (s *PasswordlessService) VerifyCode(ctx context.Context, email string, code string) (*AuthResponse, error) {
	email = normalizeEmail(email)
	code = normalizePasswordlessCode(code)

	if err := validatePasswordlessEmail(email); err != nil {
		return nil, err
	}

	if len(code) != passwordlessCodeLength {
		return nil, errors.New("invalid login code")
	}

	now := time.Now().UTC()

	var loginToken PasswordlessToken

	err := s.tokens.FindOne(
		ctx,
		bson.M{
			"email":     email,
			"purpose":   passwordlessPurposeLogin,
			"used":      false,
			"expiresAt": bson.M{"$gt": now},
		},
		options.FindOne().SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		}),
	).Decode(&loginToken)

	if err != nil {
		return nil, errors.New("invalid or expired login code")
	}

	if loginToken.Attempts >= s.maxAttempt {
		return nil, errors.New("too many attempts")
	}

	if bcrypt.CompareHashAndPassword([]byte(loginToken.CodeHash), []byte(code)) != nil {
		_, _ = s.tokens.UpdateOne(
			ctx,
			bson.M{"_id": loginToken.ID},
			bson.M{
				"$inc": bson.M{"attempts": 1},
				"$set": bson.M{"updatedAt": now},
			},
		)

		return nil, errors.New("invalid login code")
	}

	_, _ = s.tokens.UpdateOne(
		ctx,
		bson.M{"_id": loginToken.ID},
		bson.M{
			"$set": bson.M{
				"used":      true,
				"usedAt":    now,
				"updatedAt": now,
			},
		},
	)

	user, err := s.findOrCreateUser(ctx, email)
	if err != nil {
		return nil, err
	}

	userEmail := primaryUserEmail(user)
	if userEmail == "" {
		userEmail = email
	}

	token, err := s.jwtService.Generate(user.ID, userEmail, user.Role, user.SchoolID)
	if err != nil {
		return nil, errors.New("failed to create auth token")
	}

	return &AuthResponse{
		Token: token,
		User: AuthUserResponse{
			ID:       user.ID,
			Email:    userEmail,
			Role:     user.Role,
			SchoolID: user.SchoolID,
		},
	}, nil
}

func (s *PasswordlessService) findOrCreateUser(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to load user")
	}

	if user != nil {
		if !user.IsActive {
			return nil, errors.New("account is disabled")
		}

		return user, nil
	}

	now := time.Now().UTC()

	user = &User{
		ID: uuid.NewString(),
		Emails: []EmailRecord{
			{
				Address:  normalizeEmail(email),
				Verified: true,
			},
		},
		PasswordHash: "",
		Role:         "school-owner",
		Roles:        []string{"school-owner"},
		SchoolID:     "",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := s.repo.Users.InsertOne(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}
