package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"emailchatgpt/internal/service"
	"emailchatgpt/internal/utils"
)

type OTPHandler struct {
	service *service.OTPService
}

func NewOTPHandler(service *service.OTPService) *OTPHandler {
	return &OTPHandler{service: service}
}

type lookupRequest struct {
	Email string `json:"email"`
}

func (h *OTPHandler) LookupOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
		return
	}

	var req lookupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	result, err := h.service.LookupOTP(r.Context(), req.Email)
	if err != nil {
		log.Printf("otp lookup failed for %s: %v", req.Email, err)
		switch err {
		case service.ErrInvalidEmail:
			utils.WriteJSON(w, http.StatusBadRequest, result)
		case service.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, result)
		case service.ErrExpired:
			utils.WriteJSON(w, http.StatusGone, result)
		default:
			utils.WriteError(w, http.StatusInternalServerError, "gagal memproses request")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}
