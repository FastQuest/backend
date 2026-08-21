package user

import (
	"encoding/json"
	"flashquest/internal/auth"
	"flashquest/pkg/apiresp"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Get the profile data of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users/me [get]
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}
	userID := userIDValue.(uint)

	user, err := h.service.GetCurrentUser(userID)
	if err != nil {
		apiresp.WriteError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Error fetching user data")
		return
	}

	if user == nil {
		apiresp.WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": user,
	})
}
