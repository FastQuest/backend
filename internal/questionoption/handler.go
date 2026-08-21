package questionoption

import (
	"encoding/json"
	"errors"
	"flashquest/pkg/models"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

// @Summary Create question options for a question
// @Description Creates one or more question options linked to the question identified by its ID.
// @Tags QuestionOptions
// @Param id path string true "Question ID"
// @Param questionoptions body []models.QuestionOption true "List of question options to create"
// @Success 201 {object} map[string]interface{} "Question options created successfully"
// @Failure 400 {string} string "Question ID is required or invalid request body"
// @Failure 404 {string} string "Question not found"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/question-options [post]
func (h *Handler) PostQuestionOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID := vars["id"]

	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	db := h.repository.DB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	question, err := findQuestionByID(db, questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error checking question", http.StatusInternalServerError)
		}
		return
	}

	var questionoptions []models.QuestionOption
	if err := json.NewDecoder(r.Body).Decode(&questionoptions); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, questionoption := range questionoptions {
		if questionoption.Text == "" {
			http.Error(w, fmt.Sprintf("Question option text is required (index %d)", i), http.StatusBadRequest)
			return
		}
		questionoptions[i].QuestionID = question.ID
	}

	rowsAffected, err := createQuestionOptions(db, &questionoptions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving question options: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No question options were created", http.StatusInternalServerError)
		return
	}

	createdIDs := make([]uint, len(questionoptions))
	for i, questionoption := range questionoptions {
		createdIDs[i] = questionoption.ID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question options created successfully",
		"count":   rowsAffected,
		"ids":     createdIDs,
	})
}

// @Summary List question options for a question
// @Description Returns all question options associated with the specified question.
// @Tags QuestionOptions
// @Param id path string true "Question ID"
// @Success 200 {array} models.QuestionOption
// @Failure 400 {string} string "Question ID is required"
// @Failure 404 {string} string "No question options found for this question"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/question-options [get]
func (h *Handler) GetQuestionOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID := vars["id"]

	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	db := h.repository.DB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	questionoptions, err := findQuestionOptionsByQuestionID(db, questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "No question options found for this question", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question options: %v", err), http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("Found %d question options for question %s\n", len(questionoptions), questionID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questionoptions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary Fetch question options by IDs
// @Description Retrieves question options matching the provided list of question option IDs.
// @Tags QuestionOptions
// @Param request body QuestionOptionsBody true "Question option IDs payload"
// @Success 200 {array} models.QuestionOption
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /question-options/by-ids [post]
func (h *Handler) GetQuestionOptionsByIDArray(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var questionoptionsBody QuestionOptionsBody
	errConvert := json.Unmarshal(body, &questionoptionsBody)
	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	db := h.repository.DB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	questionoptions, _ := ReadQuestionOptionsByIDArray(db, questionoptionsBody.QuestionOptionIDs)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questionoptions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
