package tenant

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ==================== SEND VERIFICATION ====================

func (s *TenantService) sendVerification(ctx context.Context, t *Tenant) error {
	if t == nil {
		return errors.New("tenant is nil")
	}

	// Get recipient email
	to := strings.ToLower(strings.TrimSpace(t.Owner.Email))
	if to == "" {
		to = strings.ToLower(strings.TrimSpace(t.Email))
	}
	if to == "" {
		return errors.New("owner email missing")
	}

	if s.email.PostmarkToken == "" {
		return errors.New("postmark token missing")
	}

	from := strings.TrimSpace(s.email.FromEmail)
	if from == "" {
		from = "Blomiva <noreply@blomiva.com>"
	}

	appURL := strings.TrimRight(strings.TrimSpace(s.email.AppURL), "/")
	if appURL == "" {
		appURL = "http://localhost:8082"
	}

	verifyURL := s.buildVerificationURL(t.TenantID, to, appURL)

	ownerName := strings.TrimSpace(t.Owner.Name)
	if ownerName == "" {
		ownerName = "there"
	}

	schoolName := strings.TrimSpace(t.Name)
	if schoolName == "" {
		schoolName = "your school"
	}

	htmlBody := fmt.Sprintf(`
<!doctype html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1.0">
  <title>Verify your school</title>
</head>
<body style="margin:0;padding:0;background:#f3f4f6;font-family:Arial,Helvetica,sans-serif;color:#111827;">
  <div style="max-width:640px;margin:0 auto;padding:32px 16px;">
    <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:18px;overflow:hidden;">
      <div style="background:#1E2A44;padding:26px;text-align:center;">
        <div style="font-size:12px;letter-spacing:1.8px;text-transform:uppercase;color:#cbd5e1;font-weight:700;">Blomiva School</div>
        <h1 style="margin:8px 0 0;color:#ffffff;font-size:22px;line-height:30px;font-weight:800;">Verify your school account</h1>
      </div>

      <div style="padding:30px 26px;">
        <p style="margin:0 0 16px;font-size:15px;line-height:24px;color:#374151;">
          Hello <strong style="color:#111827;">%s</strong>,
        </p>
        <p style="margin:0 0 16px;font-size:15px;line-height:24px;color:#374151;">
          Your school <strong style="color:#111827;">%s</strong> has been created successfully on Blomiva.
        </p>
        <p style="margin:0 0 26px;font-size:15px;line-height:24px;color:#374151;">
          Verify your email address to activate the school account and continue setup.
        </p>

        <div style="text-align:center;margin:30px 0;">
          <a href="%s" style="display:inline-block;background:#1E2A44;color:#ffffff;text-decoration:none;border-radius:999px;padding:13px 24px;font-size:14px;font-weight:800;">
            Verify School Account
          </a>
        </div>

        <p style="margin:0 0 8px;font-size:12px;line-height:20px;color:#6b7280;">
          If the button does not work, copy and paste this link:
        </p>
        <p style="margin:0;font-size:12px;line-height:20px;color:#374151;word-break:break-all;">
          %s
        </p>

        <hr style="border:none;border-top:1px solid #e5e7eb;margin:26px 0;" />

        <p style="margin:0;font-size:12px;line-height:20px;color:#6b7280;">
          If you did not create this school, you can safely ignore this email.
        </p>
      </div>

      <div style="background:#f9fafb;padding:16px;text-align:center;font-size:11px;line-height:18px;color:#9ca3af;">
        © %d Blomiva. All rights reserved.
      </div>
    </div>
  </div>
</body>
</html>
`,
		html.EscapeString(ownerName),
		html.EscapeString(schoolName),
		html.EscapeString(verifyURL),
		html.EscapeString(verifyURL),
		time.Now().Year(),
	)

	textBody := fmt.Sprintf(
		`Blomiva School Verification

Hello %s,

Your school "%s" has been created successfully.

Verify here: %s

If you did not create this school, ignore this email.`,
		ownerName, schoolName, verifyURL,
	)

	payload := map[string]any{
		"From":          from,
		"To":            to,
		"Subject":       "Verify your Blomiva School Account",
		"HtmlBody":      htmlBody,
		"TextBody":      textBody,
		"MessageStream": "outbound",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.postmarkapp.com/email", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.email.PostmarkToken)

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("postmark returned status %s", res.Status)
	}

	return nil
}

// ==================== BUILD URL (DEVELOPMENT OPTIMIZED) ====================

func (s *TenantService) buildVerificationURL(tenantID, email, appURL string) string {
	token := s.signToken(tenantID, email)

	values := url.Values{}
	values.Set("tenantId", tenantID)
	values.Set("email", email)
	values.Set("token", token)

	// Always use Web URL during development (Best for testing)
	return fmt.Sprintf("%s/verify-school?%s", appURL, values.Encode())
}

// ==================== TOKEN SIGNING ====================

func (s *TenantService) signToken(id, email string) string {
	secret := strings.TrimSpace(s.email.Secret)
	if secret == "" {
		secret = "dev-secret-change-me-in-production"
	}

	msg := strings.TrimSpace(id) + "|" + strings.ToLower(strings.TrimSpace(email))

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))

	return hex.EncodeToString(h.Sum(nil))
}
