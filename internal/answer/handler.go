package answer

import (
	"encoding/json"
	"flashquest/internal/auth"
	"flashquest/internal/platform/database"
	"net/http"
)

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
func GetSubjectPerfomanceHandler(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userIDValue.(uint)
	if !ok {
		http.Error(w, "Invalid User ID format in context", http.StatusInternalServerError)
		return
	}

	userID := int(userIDUint)

	db := database.GetDB()

	subjectPerformance, err := GetUserPerfomace(db, userID)
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
func GetUserOverallPerfomanceHandler(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userIDValue.(uint)
	if !ok {
		http.Error(w, "Invalid User ID format in context", http.StatusInternalServerError)
		return
	}

	userID := int(userIDUint)

	db := database.GetDB()

	overallPerformance, err := GetUserGeralPerfomace(db, userID)
	if err != nil {
		http.Error(w, "Error fetching overall performance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(overallPerformance); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
