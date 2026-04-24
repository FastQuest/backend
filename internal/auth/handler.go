package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"flashquest/pkg/apiresp"
)

func RegisterHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apiresp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
			return
		}

		resp, err := service.Register(req)
		if err != nil {
			if errors.Is(err, ErrDuplicatedEmail) {
				apiresp.WriteError(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already in use")
				return
			}
			if errors.Is(err, ErrRoleDomainNotAllowed) {
				apiresp.WriteError(w, http.StatusUnprocessableEntity, "ROLE_DOMAIN_NOT_ALLOWED", "Email domain not allowed")
				return
			}
			apiresp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}

		apiresp.WriteJSON(w, http.StatusOK, resp)
	}
}

func LoginHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apiresp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
			return
		}

		resp, err := service.Login(req)
		if err != nil {
			if errors.Is(err, ErrInvalidCredentials) {
				apiresp.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
				return
			}
			apiresp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}

		apiresp.WriteJSON(w, http.StatusOK, resp)
	}
}
