package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"emailchatgpt/internal/config"
	"emailchatgpt/internal/handler"
	"emailchatgpt/internal/middleware"
	"emailchatgpt/internal/repository"
	"emailchatgpt/internal/service"
	"emailchatgpt/routes"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")
	cfg := config.Load()
	repo := repository.NewYahooRepository(repository.YahooConfig{
		Host:         cfg.YahooIMAPHost,
		Port:         cfg.YahooIMAPPort,
		Username:     cfg.YahooUsername,
		AppPassword:  cfg.YahooAppPassword,
		AllowedFolders: cfg.AllowedFolders,
		MaxScan:      cfg.YahooMaxScan,
	})
	otpService := service.NewOTPService(repo, cfg.AllowedFolders, cfg.MaxOTPAge)
	h := handler.NewOTPHandler(otpService)
	health := handler.NewHealthHandler()

	mux := routes.NewRouter(cfg, h, health)
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           middleware.Recover(middleware.Logging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
