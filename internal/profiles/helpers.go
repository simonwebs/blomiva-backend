package profiles

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,79}$`)
	slugRegex      = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

func CleanText(value string) string {
	return strings.TrimSpace(value)
}

func NormalizeUsername(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9_-]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-_")

	if len(value) < 3 {
		value = "user-" + value
	}

	if len(value) > 80 {
		value = value[:80]
	}

	return value
}

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")

	if value == "" {
		value = "user"
	}

	if len(value) > 100 {
		value = value[:100]
	}

	return value
}

func ValidateEmail(email string) error {
	if email == "" {
		return nil
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email address")
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return nil
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("username must be 3-80 characters and contain letters, numbers, underscore, or hyphen")
	}

	return nil
}

func SanitizeCustomKey(key string) (string, error) {
	key = strings.TrimSpace(key)

	if key == "" {
		return "", errors.New("custom key is required")
	}

	key = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(key, "")

	if key == "" || strings.Contains(key, ".") || strings.Contains(key, "$") {
		return "", errors.New("invalid custom key")
	}

	return key, nil
}

func GetUserEmail(user *User) string {
	if user == nil || len(user.Emails) == 0 {
		return ""
	}
	return strings.TrimSpace(strings.ToLower(user.Emails[0].Address))
}