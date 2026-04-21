package service

import (
	"context"
	"testing"
	"time"

	"emailchatgpt/internal/model"
	"emailchatgpt/internal/repository"
)

type stubRepository struct {
	messages []model.EmailMessage
}

func (s stubRepository) ListMessages(_ context.Context, _ string) ([]model.EmailMessage, error) {
	return append([]model.EmailMessage(nil), s.messages...), nil
}

var _ repository.Repository = stubRepository{}

func TestLookupOTPAcceptsLocalizedChatGPTSubject(t *testing.T) {
	now := time.Now()
	repo := stubRepository{messages: []model.EmailMessage{{
		ID:         "1",
		Folder:     "INBOX",
		Sender:     "noreply@chatgpt.com",
		Recipient:  "user@example.com",
		Subject:    "Kode masuk ChatGPT",
		Text:       "Kode verifikasi sementara: 123456",
		ReceivedAt: now.Add(-2 * time.Minute),
	}}}

	service := NewOTPService(repo, []string{"INBOX"}, 5*time.Minute)
	result, err := service.LookupOTP(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("LookupOTP returned error: %v", err)
	}

	if result.Status != "found" {
		t.Fatalf("expected status found, got %q", result.Status)
	}
	if result.OTP != "123456" {
		t.Fatalf("expected OTP 123456, got %q", result.OTP)
	}
}

func TestLookupOTPUsesBodyHintsWhenSubjectIsGeneric(t *testing.T) {
	now := time.Now()
	repo := stubRepository{messages: []model.EmailMessage{{
		ID:         "1",
		Folder:     "OTP",
		Sender:     "noreply@chatgpt.com",
		Recipient:  "user@example.com",
		Subject:    "Security alert",
		Text:       "Enter this temporary verification code to continue: 654321",
		ReceivedAt: now.Add(-1 * time.Minute),
	}}}

	service := NewOTPService(repo, []string{"OTP"}, 5*time.Minute)
	result, err := service.LookupOTP(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("LookupOTP returned error: %v", err)
	}

	if result.Status != "found" {
		t.Fatalf("expected status found, got %q", result.Status)
	}
	if result.OTP != "654321" {
		t.Fatalf("expected OTP 654321, got %q", result.OTP)
	}
}
