package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"regexp"
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

var chatGPTSubjectHints = []string{
	"your temporary chatgpt login code",
	"chatgpt login code",
	"temporary chatgpt login code",
	"kode masuk chatgpt",
	"kode login chatgpt",
}

var chatGPTBodyHints = []string{
	"enter this temporary verification code to continue",
	"temporary verification code",
	"kode verifikasi sementara",
	"kode verifikasi ini",
	"gunakan kode verifikasi ini",
}

var otpPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?is)Enter this temporary verification code to continue:\s*([0-9]{4,8})`),
	regexp.MustCompile(`(?is)Kode verifikasi(?: Anda)?(?: sementara)?[:\s]+([0-9]{4,8})`),
	regexp.MustCompile(`(?is)Gunakan kode(?: verifikasi)?(?: ini)?[:\s]+([0-9]{4,8})`),
}

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
		if !isChatGPTMessage(message.Subject, message.Body, message.Text) {
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

		otp := extractChatGPTOTP(message)

		return model.OTPResult{
			Status:     "found",
			Message:    "email ditemukan",
			Email:      email,
			OTP:        otp,
			Subject:    message.Subject,
			Sender:     message.Sender,
			Folder:     message.Folder,
			HTML:       message.Body,
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

func extractChatGPTOTP(message model.EmailMessage) string {
	text := strings.TrimSpace(message.Text)
	if text == "" {
		text = stripHTML(message.Body)
	}
	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	for _, pattern := range otpPatterns {
		if match := pattern.FindStringSubmatch(text); len(match) >= 2 {
			return match[1]
		}
	}
	return ""
}

func isChatGPTMessage(subject, body, text string) bool {
	joined := strings.ToLower(strings.TrimSpace(subject) + " " + strings.TrimSpace(body) + " " + strings.TrimSpace(text))
	for _, hint := range chatGPTSubjectHints {
		if strings.Contains(joined, hint) {
			return true
		}
	}
	for _, hint := range chatGPTBodyHints {
		if strings.Contains(joined, hint) {
			return true
		}
	}
	return false
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
