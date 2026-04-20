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
