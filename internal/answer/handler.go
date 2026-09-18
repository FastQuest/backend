package answer

import (
	"encoding/json"
	"net/http"

	"flashquest/internal/auth"
)

// Handler is the HTTP layer of the answer domain. It depends only on the
// Service interface, never on the database.
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetSubjectPerfomance godoc
// @Summary Get user subjections performance
// @Description Get all answers performance from authenticated user, optionally filtered by question set
// @Tags Answers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved answers"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /answers/performance [get]
func (h *Handler) GetSubjectPerfomanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	subjectPerformance, err := h.service.GetSubjectPerformance(userID)
	if err != nil {
		http.Error(w, "Error fetching subject performance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subjectPerformance); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetUserOverallPerfomance godoc
// @Summary Get user overall performance
// @Description Get user's overall performance metrics from authenticated user
// @Tags Answers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved answers"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /answers/overall-performance [get]
func (h *Handler) GetUserOverallPerfomanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(w, r)
	if !ok {
		return
	}

	overallPerformance, err := h.service.GetOverallPerformance(userID)
	if err != nil {
		http.Error(w, "Error fetching overall performance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(overallPerformance); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// userIDFromContext reads the authenticated user id, writing the response
// itself when it is missing or malformed.
func userIDFromContext(w http.ResponseWriter, r *http.Request) (int, bool) {
	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return 0, false
	}

	userIDUint, ok := userIDValue.(uint)
	if !ok {
		http.Error(w, "Invalid User ID format in context", http.StatusInternalServerError)
		return 0, false
	}

	return int(userIDUint), true
}
