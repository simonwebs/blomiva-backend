package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

type EmailSender interface {
	Send(ctx context.Context, to string, subject string, text string) error
}

type PostmarkEmailSender struct {
	token string
	from  string
}

func NewPostmarkEmailSender() *PostmarkEmailSender {
	return &PostmarkEmailSender{
		token: envFirst("POSTMARK_SERVER_TOKEN", "POSTMARK_API_KEY", "POSTMARK_TOKEN"),
		from:  envDefault("POSTMARK_FROM_EMAIL", "Blomiva <noreply@blomiva.com>"),
	}
}

func (s *PostmarkEmailSender) Send(ctx context.Context, to string, subject string, text string) error {
	if strings.TrimSpace(s.token) == "" || strings.TrimSpace(to) == "" {
		return nil
	}

	payload := map[string]string{
		"From":     s.from,
		"To":       to,
		"Subject":  subject,
		"TextBody": text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.postmarkapp.com/email",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.token)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func envFirst(keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return ""
}

func envDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
