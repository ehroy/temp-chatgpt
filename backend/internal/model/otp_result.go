package model

import "time"

type OTPResult struct {
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	Email      string    `json:"email"`
	OTP        string    `json:"otp,omitempty"`
	Subject    string    `json:"subject,omitempty"`
	Sender     string    `json:"sender,omitempty"`
	Folder     string    `json:"folder,omitempty"`
	HTML       string    `json:"html,omitempty"`
	ReceivedAt time.Time `json:"receivedAt,omitempty"`
	ExpiresAt  time.Time `json:"expiresAt,omitempty"`
}
