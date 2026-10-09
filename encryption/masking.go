package encryption

import (
	"strings"
)

// MaskEmail masks an email (e.g. "john.doe@example.com" -> "j***e@example.com").
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" {
		return email
	}
	user := parts[0]
	domain := parts[1]

	if len(user) <= 2 {
		return user[:1] + "***@" + domain
	}
	return string(user[0]) + "***" + string(user[len(user)-1]) + "@" + domain
}

// MaskPhone masks a phone number (e.g. "+6281234567890" -> "+62812****7890").
func MaskPhone(phone string) string {
	clean := strings.TrimSpace(phone)
	if len(clean) <= 6 {
		return clean
	}
	prefixLen := 4
	if strings.HasPrefix(clean, "+") {
		prefixLen = 5
	}
	if len(clean) <= prefixLen+4 {
		return clean[:2] + "****" + clean[len(clean)-2:]
	}
	return clean[:prefixLen] + "****" + clean[len(clean)-4:]
}
