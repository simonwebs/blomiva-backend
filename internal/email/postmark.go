package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PostmarkSender struct {
	APIKey string
	From   string
}

func NewPostmarkSender(apiKey, from string) *PostmarkSender {
	return &PostmarkSender{
		APIKey: apiKey,
		From:   from,
	}
}

func (p *PostmarkSender) Send(to string, subject string, html string) error {
	payload := map[string]any{
		"From":     p.From,
		"To":       to,
		"Subject":  subject,
		"HtmlBody": html,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://api.postmarkapp.com/email",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", p.APIKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("postmark error: %s", res.Status)
	}

	return nil
}
