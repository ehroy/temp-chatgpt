package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr            string
	AuthToken       string
	AllowedFolders  []string
	MaxOTPAge       time.Duration
	YahooIMAPHost   string
	YahooIMAPPort   string
	YahooUsername   string
	YahooAppPassword string
	YahooMaxScan    int
}

func Load() Config {
	return Config{
		Addr:             envOrDefault("APP_ADDR", "127.0.0.1:9001"),
		AuthToken:        envOrDefault("APP_AUTH_TOKEN", "dev-token"),
		AllowedFolders:   parseCSV(envOrDefault("YAHOO_ALLOWED_FOLDERS", "INBOX,OTP")),
		MaxOTPAge:        parseDurationMinutes(envOrDefault("OTP_MAX_AGE_MINUTES", "5")),
		YahooIMAPHost:    envOrDefault("YAHOO_IMAP_HOST", "imap.mail.yahoo.com"),
		YahooIMAPPort:    envOrDefault("YAHOO_IMAP_PORT", "993"),
		YahooUsername:    envOrDefault("YAHOO_EMAIL", ""),
		YahooAppPassword: envOrDefault("YAHOO_APP_PASSWORD", ""),
		YahooMaxScan:     parseInt(envOrDefault("YAHOO_MAX_SCAN", "25"), 25),
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	folders := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			folders = append(folders, trimmed)
		}
	}
	return folders
}

func parseDurationMinutes(value string) time.Duration {
	minutes, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || minutes <= 0 {
		minutes = 5
	}
	return time.Duration(minutes) * time.Minute
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
