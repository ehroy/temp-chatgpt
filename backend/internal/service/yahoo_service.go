package service

// YahooService is reserved for the real mailbox adapter.
// The current scaffold keeps the OTP logic behind repository/service boundaries
// so the Yahoo IMAP implementation can be swapped in without changing handlers.
type YahooService struct{}
