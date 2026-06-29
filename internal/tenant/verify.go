package tenant

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *TenantService) VerifySchool(ctx context.Context, tenantID, email, token string) (*Tenant, error) {
	t, err := s.repo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	expected := s.sign(tenantID, email)

	if !hmac.Equal([]byte(expected), []byte(token)) {
		return nil, errors.New("invalid token")
	}

	now := time.Now().UTC()

	update := bson.M{
		"verification.emailVerified": true,
		"verification.verifiedAt":    now,
		"status":                     StatusActive,
		"updatedAt":                  now,
	}

	return s.repo.Update(ctx, t.ID, update)
}

func (s *TenantService) sign(id, email string) string {
	secret := s.email.Secret
	if secret == "" {
		secret = "dev"
	}

	m := id + "|" + strings.ToLower(email)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(m))

	return hex.EncodeToString(h.Sum(nil))
}
