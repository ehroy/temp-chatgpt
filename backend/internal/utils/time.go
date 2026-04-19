package utils

import (
	"strings"
	"time"
)

func IsToday(candidate, now time.Time) bool {
	y1, m1, d1 := candidate.In(now.Location()).Date()
	y2, m2, d2 := now.In(now.Location()).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func MaskOTP(otp string) string {
	if len(otp) <= 2 {
		return "**"
	}
	return otp[:2] + strings.Repeat("*", len(otp)-2)
}
