package submission

import (
	"encoding/json"
	"flashquest/internal/auth"
	"fmt"
	"math"
	"net/http"
	"strconv"
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

// GetUserSubmissions godoc
// @Summary Get user submissions
// @Description Get all submissions from authenticated user, optionally filtered by question set
// @Tags Submissions
// @Accept json
// @Produce json
// @Param question_set_id query integer false "Filter by question set ID"
// @Param page query integer false "Page number (default: 1)"
// @Param perPage query integer false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Successfully retrieved submissions"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users/submissions [get]
func GetUserSubmissions(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.ContextKeyUserID)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	// Parse query parameters
	queryParams := r.URL.Query()

	var questionSetID *uint
	if qsID := queryParams.Get("question_set_id"); qsID != "" {
		id, err := strconv.ParseUint(qsID, 10, 32)
		if err != nil {
			http.Error(w, "Invalid question_set_id", http.StatusBadRequest)
			return
		}
		qsIDUint := uint(id)
		questionSetID = &qsIDUint
	}

	// Pagination
	page := 1
	pageSize := 10
	if p := queryParams.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := queryParams.Get("perPage"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	db := getDB()
	submissions, total, err := getSubmissionsByUserID(db, userID, questionSetID, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching submissions: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": submissions,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     pageSize,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(pageSize))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
