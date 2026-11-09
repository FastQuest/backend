package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Handler is the HTTP layer of the question domain. It depends only on the
// Service interface, never on the database.
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateQuestion godoc
// @Summary Create a new question
// @Description Creates a single question from the provided JSON object.
// @Tags Questions
// @Accept json
// @Produce json
// @Param question body CreateRequest true "Question to create"
// @Success 201 {object} CreateResponse "Successfully created question"
// @Failure 400 {string} string "Invalid request body format or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /questions [post]
func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var request CreateRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request body format: expected a single question object", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateQuestion(request)
	if err != nil {
		writeCreateError(w, err)
		return
	}

	sendJSON(w, created, http.StatusCreated)
}

// CreateQuestionsBatch godoc
// @Summary Create multiple questions
// @Description Creates a batch of questions from the provided JSON array. Questions are created in order and the ones before an invalid or failing question stay created.
// @Tags Questions
// @Accept json
// @Produce json
// @Param questions body []CreateRequest true "Questions to create"
// @Success 201 {array} CreateResponse "Successfully created questions"
// @Failure 400 {string} string "Invalid request body format or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/batch [post]
func (h *Handler) CreateQuestionsBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var requests []CreateRequest
	if err := json.Unmarshal(body, &requests); err != nil || len(requests) == 0 {
		http.Error(w, "Invalid request body format: expected a non-empty array of question objects", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateQuestions(requests)
	if err != nil {
		writeCreateError(w, err)
		return
	}

	sendJSON(w, created, http.StatusCreated)
}

// GetQuestions godoc
// @Summary Get paginated questions
// @Description Retrieve a paginated list of questions with optional sorting, filtering, and eager loading of relationships.
// @Tags Questions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param perPage query int false "Items per page (max 100)" default(10)
// @Param orderBy query string false "Order by field (e.g., created_at desc)"
// @Param source query int false "Filter by source ID"
// @Param source_exam_instance_id query int false "Filter by source exam instance ID"
// @Param include query string false "Comma-separated list of relations to include (source, question-options, user and subject)"
// @Param statement query string false "Filter by statement (case insensitive)"
// @Param question_option query string false "Filter by question option text (case insensitive)"
// @Param subject query string false "Filter by subject ID"
// @Param topic query string false "Filter by topic"
// @Param source query string false "Filter by source ID"
// @Param year query int false "Filter by exam year or creation year"
// @Param list query bool false "Filter questions that belong (true) or not (false) to a list"
// @Success 200 {object} map[string]interface{} "Paginated response containing data and pagination metadata"
// @Failure 500 {string} string "Internal server error"
// @Router /questions [get]
func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := parseInt(query.Get("page"), 1)
	limit := parseInt(query.Get("perPage"), 10)
	if limit > 100 {
		limit = 100
	}

	orderBy := query.Get("orderBy")
	if orderBy == "" {
		orderBy = "created_at desc"
	}

	filter := parseFilter(query)

	total, err := h.service.CountQuestions(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting questions: %v", err), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	questions, err := h.service.GetQuestions(filter, orderBy, parseIncludes(query), offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", err), http.StatusInternalServerError)
		return
	}

	sendPaginatedResponse(w, questions, total, limit, page)
}

// GetQuestion godoc
// @Summary Get a specific question
// @Description Retrieve a specific question by its ID.
// @Tags Questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param include query string false "Comma-separated list of relations to include"
// @Success 200 {object} Response
// @Failure 400 {string} string "ID parameter is required"
// @Failure 404 {string} string "Question not found"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id} [get]
func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	question, err := h.service.GetQuestion(id, parseIncludes(r.URL.Query()))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetQuestionFilters godoc
// @Summary Get question filters
// @Description Retrieve subjects, sources, source exam instances, and years available for question filters.
// @Tags Questions
// @Accept json
// @Produce json
// @Success 200 {object} FiltersResponse
// @Failure 500 {string} string "Internal server error"
// @Router /questions/filters [get]
func (h *Handler) GetQuestionFiltersHandler(w http.ResponseWriter, r *http.Request) {
	filters, err := h.service.GetQuestionFilters()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching question filters: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filters); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetQuestionsByArray godoc
// @Summary Get multiple questions by array of IDs
// @Description Retrieve multiple questions by providing an array of their IDs in the request body.
// @Tags Questions
// @Accept json
// @Produce json
// @Param include query string false "Comma-separated list of relations to include"
// @Param request body IDsRequest true "Array of question IDs"
// @Success 200 {array} Response
// @Failure 400 {string} string "Invalid JSON body or empty IDs array"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/by-ids [post]
func (h *Handler) GetQuestionsByArray(w http.ResponseWriter, r *http.Request) {
	var req IDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, "Invalid JSON body or empty IDs array", http.StatusBadRequest)
		return
	}

	questions, err := h.service.GetQuestionsByIDs(req.IDs, parseIncludes(r.URL.Query()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSON(w, questions, http.StatusOK)
}

// DeleteQuestion godoc
// @Summary Delete a question
// @Description Delete a specific question by its ID.
// @Tags Questions
// @Param id path string true "Question ID"
// @Success 200 {string} string "Successfully deleted"
// @Failure 400 {string} string "ID parameter is required"
// @Failure 404 {string} string "Question not found"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id} [delete]
func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to delete question ID: %s", id)
	rows, err := h.service.DeleteQuestion(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting question: %v", err), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	log.Printf("Question ID %s deleted successfully", id)
}

// writeCreateError maps the creation errors: validation failures are the
// client's fault, anything else comes from persistence.
func writeCreateError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrMissingRequiredFields) || errors.Is(err, ErrBatchMissingRequiredFields) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, fmt.Sprintf("Error creating question: %v", err), http.StatusInternalServerError)
}

func parseFilter(query url.Values) QuestionFilterDTO {
	return QuestionFilterDTO{
		Statement:      query.Get("statement"),
		QuestionOption: query.Get("question_option"),
		Subject:        query.Get("subject"),
		Topic:          query.Get("topic"),
		Source:         query.Get("source"),
		Year:           query.Get("year"),
		List:           query.Get("list"),
	}
}

func parseIncludes(query url.Values) []string {
	includeParam := query.Get("include")
	if includeParam == "" {
		return nil
	}
	return strings.Split(includeParam, ",")
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func sendPaginatedResponse(w http.ResponseWriter, data interface{}, total int64, limit, page int) {
	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}
	sendJSON(w, response, http.StatusOK)
}

func parseInt(value string, defaultValue int) int {
	if result, err := strconv.Atoi(value); err == nil && result > 0 {
		return result
	}
	return defaultValue
}
