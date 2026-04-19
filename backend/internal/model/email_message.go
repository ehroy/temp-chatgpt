package model

import "time"

type EmailMessage struct {
	ID         string    `json:"id"`
	Folder     string    `json:"folder"`
	Sender     string    `json:"sender"`
	Recipient  string    `json:"recipient"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	Text       string    `json:"text,omitempty"`
	ReceivedAt time.Time `json:"receivedAt"`
}
