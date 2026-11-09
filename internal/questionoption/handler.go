package questionoption

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary Create question options for a question
// @Description Creates one or more question options linked to the question identified by its ID.
// @Tags QuestionOptions
// @Param id path string true "Question ID"
// @Param questionoptions body []CreateRequest true "List of question options to create"
// @Success 201 {object} map[string]interface{} "Question options created successfully"
// @Failure 400 {string} string "Question ID is required or invalid request body"
// @Failure 404 {string} string "Question not found"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/question-options [post]
func (h *Handler) PostQuestionOptions(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]
	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	var requests []CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateQuestionOptions(questionID, requests)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, ErrQuestionNotFound):
			http.Error(w, "Question not found", http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Error saving question options: %v", err), http.StatusInternalServerError)
		}
		return
	}
	if result.Count == 0 {
		http.Error(w, "No question options were created", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question options created successfully",
		"count":   result.Count,
		"ids":     result.IDs,
	}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary List question options for a question
// @Description Returns all question options associated with the specified question.
// @Tags QuestionOptions
// @Param id path string true "Question ID"
// @Success 200 {array} Response
// @Failure 400 {string} string "Question ID is required"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/question-options [get]
func (h *Handler) GetQuestionOptions(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]
	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	options, err := h.service.GetQuestionOptionsByQuestionID(questionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching question options: %v", err), http.StatusInternalServerError)
		return
	}
	sendJSON(w, options, http.StatusOK)
}

// @Summary Fetch question options by IDs
// @Description Retrieves question options matching the provided list of question option IDs.
// @Tags QuestionOptions
// @Param request body QuestionOptionsBody true "Question option IDs payload"
// @Success 200 {array} Response
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /question-options/by-ids [post]
func (h *Handler) GetQuestionOptionsByIDArray(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var request QuestionOptionsBody
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	options, err := h.service.GetQuestionOptionsByIDs(request.QuestionOptionIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching question options: %v", err), http.StatusInternalServerError)
		return
	}
	sendJSON(w, options, http.StatusOK)
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
