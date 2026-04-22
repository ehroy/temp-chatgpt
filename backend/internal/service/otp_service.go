package service

import (
	"context"
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"mime/quotedprintable"
	"regexp"
	"unicode"
	"strings"
	"time"

	"emailchatgpt/internal/model"
	"emailchatgpt/internal/repository"
	"emailchatgpt/internal/utils"
)

var (
	ErrInvalidEmail = errors.New("email tidak valid")
	ErrNotFound     = errors.New("otp tidak ditemukan")
	ErrExpired      = errors.New("otp sudah kadaluarsa")
	ErrNotAllowed   = errors.New("folder email tidak diizinkan")
)

var otpPattern = regexp.MustCompile(`(?m)\b([0-9]{4,8})\b`)

type OTPService struct {
	repo           repository.Repository
	allowedFolders map[string]struct{}
	maxAge         time.Duration
}

func NewOTPService(repo repository.Repository, folders []string, maxAge time.Duration) *OTPService {
	allowed := make(map[string]struct{}, len(folders))
	for _, folder := range folders {
		allowed[strings.ToLower(strings.TrimSpace(folder))] = struct{}{}
	}
	return &OTPService{
		repo:           repo,
		allowedFolders: allowed,
		maxAge:         maxAge,
	}
}

func (s *OTPService) LookupOTP(ctx context.Context, email string) (model.OTPResult, error) {
	if !utils.ValidEmail(email) {
		return model.OTPResult{Status: "invalid_email", Message: ErrInvalidEmail.Error(), Email: email}, ErrInvalidEmail
	}

	messages, err := s.repo.ListMessages(ctx, email)
	if err != nil {
		return model.OTPResult{}, err
	}

	now := time.Now()
	for _, message := range messages {
		if strings.TrimSpace(message.Recipient) != "" && !strings.Contains(strings.ToLower(message.Recipient), strings.ToLower(strings.TrimSpace(email))) {
			continue
		}
		if !s.folderAllowed(message.Folder) {
			continue
		}
		if !utils.IsToday(message.ReceivedAt, now) {
			continue
		}
		if now.Sub(message.ReceivedAt) > s.maxAge {
			return model.OTPResult{Status: "expired", Message: ErrExpired.Error(), Email: email, Subject: message.Subject, Sender: message.Sender, Folder: message.Folder, ReceivedAt: message.ReceivedAt, ExpiresAt: message.ReceivedAt.Add(s.maxAge)}, ErrExpired
		}

		text := message.Text
		if strings.TrimSpace(text) == "" {
			text = stripHTML(message.Body)
		}
		text = normalizeEmailText(text)
		otp := extractOTP(text)

		return model.OTPResult{
			Status:     "found",
			Message:    "email ditemukan",
			Email:      email,
			OTP:        otp,
			Text:       text,
			Subject:    message.Subject,
			Sender:     message.Sender,
			Folder:     message.Folder,
			ReceivedAt: message.ReceivedAt,
			ExpiresAt:  message.ReceivedAt.Add(s.maxAge),
		}, nil
	}

	return model.OTPResult{Status: "not_found", Message: ErrNotFound.Error(), Email: email}, ErrNotFound
}

func (s *OTPService) folderAllowed(folder string) bool {
	if len(s.allowedFolders) == 0 {
		return true
	}
	_, ok := s.allowedFolders[strings.ToLower(strings.TrimSpace(folder))]
	return ok
}

func (s *OTPService) DebugString() string {
	return fmt.Sprintf("otp-service(maxAge=%s)", s.maxAge)
}

func extractOTP(text string) string {
	match := otpPattern.FindStringSubmatch(text)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

func normalizeEmailText(text string) string {
	text = strings.TrimSpace(html.UnescapeString(text))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "=\n", "")

	if looksQuotedPrintable(text) {
		if decoded, err := io.ReadAll(quotedprintable.NewReader(bytes.NewBufferString(text))); err == nil && len(decoded) > 0 {
			text = string(decoded)
			text = strings.ReplaceAll(text, "\r\n", "\n")
			text = strings.ReplaceAll(text, "\r", "\n")
			text = strings.TrimSpace(html.UnescapeString(text))
		}
	}

	var b strings.Builder
	lastSpace := false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
			lastSpace = false
		case unicode.IsSpace(r):
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
		default:
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
		}
	}
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(b.String(), " "))
}

func looksQuotedPrintable(text string) bool {
	return strings.Contains(text, "=") && regexp.MustCompile(`=[0-9A-Fa-f]{2}`).MatchString(text)
}

func stripHTML(value string) string {
	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n",
		"</div>", "\n",
		"</tr>", "\n",
		"</li>", "\n",
	)
	cleaned := replacer.Replace(value)
	cleaned = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`(?is)<[^>]+>`).ReplaceAllString(cleaned, " ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	return strings.TrimSpace(cleaned)
}
