package tenant

import (
	"regexp"
	"strings"
)

func normalize(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	v = re.ReplaceAllString(v, "-")
	return strings.Trim(v, "-")
}
