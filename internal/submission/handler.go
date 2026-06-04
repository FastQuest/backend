package submission

import (
	"encoding/json"
	"flashquest/internal/auth"
	"fmt"
	"net/http"
)

// CreateSubmission godoc
// @Summary Create a new submission
// @Description Creates a new submission with user answers for a question set
// @Tags Submissions
// @Accept json
// @Produce json
// @Param submission body CreateSubmissionRequest true "Submission with user answers"
// @Success 200 {object} map[string]interface{} "Successfully created submission"
// @Failure 400 {string} string "Invalid JSON body or request format"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /submissions [post]
func CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req CreateSubmissionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		println(err.Error())
		return
	}

	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	req.UserID = userIDValue.(uint)

	submission, err := CreateSubmissionPayload(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating submission: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submission)
}
