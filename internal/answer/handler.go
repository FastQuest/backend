package answer

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

// @Summary Create answers for a question
// @Description Creates one or more answers linked to the question identified by its ID.
// @Param id path string true "Question ID"
// @Param answers body []models.Answer true "List of answers to create"
// @Success 201 {object} map[string]interface{} "Answers created successfully"
// @Failure 400 {string} string "Question ID is required or invalid request body"
// @Failure 404 {string} string "Question not found"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/answers [post]
func PostAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID := vars["id"]

	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	db := getDB()
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

	var answers []models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answers); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, answer := range answers {
		if answer.Text == "" {
			http.Error(w, fmt.Sprintf("Answer text is required (index %d)", i), http.StatusBadRequest)
			return
		}
		answers[i].QuestionID = question.ID
	}

	rowsAffected, err := createAnswers(db, &answers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving answers: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No answers were created", http.StatusInternalServerError)
		return
	}

	createdIDs := make([]uint, len(answers))
	for i, answer := range answers {
		createdIDs[i] = answer.ID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Answers created successfully",
		"count":   rowsAffected,
		"ids":     createdIDs,
	})
}

// @Summary List answers for a question
// @Description Returns all answers associated with the specified question.
// @Param id path string true "Question ID"
// @Success 200 {array} models.Answer
// @Failure 400 {string} string "Question ID is required"
// @Failure 404 {string} string "No answers found for this question"
// @Failure 500 {string} string "Internal server error"
// @Router /questions/{id}/answers [get]
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID := vars["id"]

	if questionID == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	answers, err := findAnswersByQuestionID(db, questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "No answers found for this question", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching answers: %v", err), http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("Found %d answers for question %s\n", len(answers), questionID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answers); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary Fetch answers by IDs
// @Description Retrieves answers matching the provided list of answer IDs.
// @Param request body AnswersBody true "Answer IDs payload"
// @Success 200 {array} models.Answer
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /answers/by-ids [post]
func GetAnswersByIDArray(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var answersBody AnswersBody
	errConvert := json.Unmarshal(body, &answersBody)
	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	answers, _ := readAnswersByIDArray(db, answersBody.AnswersIDs)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answers); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
