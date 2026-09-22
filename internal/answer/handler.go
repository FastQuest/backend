package answer

import (
	"encoding/json"
	"flashquest/internal/appcontext"
	"flashquest/pkg/apiresp"
	"net/http"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
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
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	subjectPerformance, err := h.repository.GetUserPerfomace(int(userID))
	if err != nil {
		apiresp.WriteError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Error fetching subject performance")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": subjectPerformance,
	})
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
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	overallPerformance, err := h.repository.GetUserGeralPerfomace(int(userID))
	if err != nil {
		apiresp.WriteError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Error fetching overall performance")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": overallPerformance,
	})
}
