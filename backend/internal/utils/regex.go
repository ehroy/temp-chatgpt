package utils

import (
	"regexp"
	"strings"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
var otpPattern = regexp.MustCompile(`\b(\d{4,8})\b`)

func ValidEmail(value string) bool {
	return emailPattern.MatchString(strings.TrimSpace(value))
}

func ExtractOTP(text string) (string, bool) {
	matches := otpPattern.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}

func IsLikelyOTPMessage(subject, body string, patterns []string) bool {
	joined := strings.ToLower(subject + " " + body)
	for _, pattern := range patterns {
		if strings.Contains(joined, strings.ToLower(pattern)) {
			return true
		}
	}
	_, ok := ExtractOTP(joined)
	return ok
}
