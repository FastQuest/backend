package submission

import (
	"encoding/json"
	"errors"
	"flashquest/internal/appcontext"
	"flashquest/pkg/apiresp"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateSubmission godoc
// @Summary Create a new submission
// @Description Creates a new submission with user answers for a question set
// @Tags Submissions
// @Accept json
// @Produce json
// @Param submission body CreateSubmissionRequest true "Submission with user answers"
// @Success 200 {object} Response "Successfully created submission"
// @Failure 400 {string} string "Invalid JSON body or request format"
// @Failure 404 {string} string "Question set not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /submissions [post]
func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req CreateSubmissionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}
	req.UserID = userID

	submission, err := h.service.CreateSubmissionPayload(req)
	if err != nil {
		if errors.Is(err, ErrQuestionSetNotFound) {
			http.Error(w, "Question set not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error creating submission: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, submission, http.StatusOK)
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
// @Success 200 {object} ListResponse "Successfully retrieved submissions"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /submissions [get]
func (h *Handler) GetUserSubmissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

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

	submissions, total, err := h.service.GetSubmissionsByUserID(userID, questionSetID, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching submissions: %v", err), http.StatusInternalServerError)
		return
	}

	response := ListResponse{
		Data: submissions,
		Pagination: PaginationResponse{
			Total:       total,
			PerPage:     pageSize,
			CurrentPage: page,
			LastPage:    int(math.Ceil(float64(total) / float64(pageSize))),
		},
	}

	writeJSON(w, response, http.StatusOK)
}

// GetSubmission godoc
// @Summary Get a specific submission
// @Description Get a specific submission with optional includes for user and answers
// @Tags Submissions
// @Accept json
// @Produce json
// @Param id path integer true "Submission ID"
// @Param include query string false "Comma-separated includes (user, answers)"
// @Success 200 {object} Response "Successfully retrieved submission"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Submission not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /submissions/{id} [get]
func (h *Handler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		apiresp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	submissionID := mux.Vars(r)["id"]
	if submissionID == "" {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "ID parameter is required")
		return
	}

	id, err := strconv.ParseUint(submissionID, 10, 32)
	if err != nil {
		apiresp.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid submission ID")
		return
	}

	includes := r.URL.Query().Get("include")

	submission, err := h.service.GetSubmissionByID(uint(id), includes)
	if err != nil {
		apiresp.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Submission not found")
		return
	}

	if submission.UserID != userID {
		apiresp.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Unauthorized")
		return
	}

	writeJSON(w, submission, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
