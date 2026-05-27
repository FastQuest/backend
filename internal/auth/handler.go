package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"flashquest/pkg/apiresp"
)

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Registers a new user based on their email domain and returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration payload"
// @Success      200  {object}  AuthResponse "Successfully registered"
// @Failure      400  {object}  map[string]interface{} "Invalid request payload or validation error"
// @Failure      409  {object}  map[string]interface{} "Email already in use"
// @Failure      422  {object}  map[string]interface{} "Email domain not allowed"
// @Router       /api/auth/register [post]
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

// LoginHandler godoc
// @Summary      Authenticate a user
// @Description  Authenticates a user with email and password, returning access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200  {object}  AuthResponse "Successfully authenticated"
// @Failure      400  {object}  map[string]interface{} "Invalid request payload or validation error"
// @Failure      401  {object}  map[string]interface{} "Invalid email or password"
// @Router       /api/auth/login [post]
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