package routes

import (
	"net/http"

	"emailchatgpt/internal/config"
	"emailchatgpt/internal/handler"
	"emailchatgpt/internal/middleware"
)

func NewRouter(cfg config.Config, otp *handler.OTPHandler, health *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.Health)
	mux.Handle("/api/otp/lookup", middleware.Auth(cfg.AuthToken)(http.HandlerFunc(otp.LookupOTP)))
	return middleware.CORS(mux)
}
