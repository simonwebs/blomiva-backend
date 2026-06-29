package auth

import (
	"fmt"
	"os"

	postmark "github.com/keighl/postmark"
)

type PostmarkEmailSender struct {
	client *postmark.Client
	from   string
}

func NewPostmarkEmailSender() *PostmarkEmailSender {
	token := os.Getenv("POSTMARK_SERVER_TOKEN")
	from := os.Getenv("POSTMARK_FROM_EMAIL")

	if token == "" {
		panic("POSTMARK_SERVER_TOKEN is required")
	}

	if from == "" {
		panic("POSTMARK_FROM_EMAIL is required")
	}

	return &PostmarkEmailSender{
		client: postmark.NewClient(token, ""),
		from:   from,
	}
}

func (s *PostmarkEmailSender) Send(to string, subject string, html string) error {
	email := postmark.Email{
		From:     s.from,
		To:       to,
		Subject:  subject,
		HtmlBody: html,
		TextBody: "Your Blomiva School login code is in this email.",
	}

	_, err := s.client.SendEmail(email)
	if err != nil {
		return fmt.Errorf("postmark send failed: %w", err)
	}

	return nil
}
