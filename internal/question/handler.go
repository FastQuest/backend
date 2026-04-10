package question

import (
	"encoding/json"
	"errors"
	filters "flashquest/pkg"
	"flashquest/pkg/models"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var questionsToProcess []QuestionInput
	var createdQuestions []models.Question

	var questionArray []QuestionInput
	errArray := json.Unmarshal(body, &questionArray)

	if errArray == nil && len(questionArray) > 0 {
		questionsToProcess = questionArray
	} else {
		var singleQuestion QuestionInput
		errSingle := json.Unmarshal(body, &singleQuestion)

		if errSingle == nil && (singleQuestion.Statement != "" || singleQuestion.UserID != 0) {
			questionsToProcess = []QuestionInput{singleQuestion}
		} else {
			http.Error(w, "Invalid request body format: expected single question object or non-empty array of objects", http.StatusBadRequest)
			return
		}
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	for _, input := range questionsToProcess {
		if input.Statement == "" || input.UserID == 0 {
			if len(questionsToProcess) > 1 {
				http.Error(w, "One or more questions are missing Statement or User ID in the batch request", http.StatusBadRequest)
				return
			}
			http.Error(w, "Statement and User ID are required", http.StatusBadRequest)
			return
		}

		question := models.Question{
			Statement:            input.Statement,
			SubjectID:            input.SubjectID,
			UserID:               input.UserID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			SourceExamInstanceID: input.SourceExamInstanceID,
		}

		if err := createQuestion(db, &question); err != nil {
			http.Error(w, fmt.Sprintf("Error creating question: %v", err), http.StatusInternalServerError)
			return
		}

		createdQuestions = append(createdQuestions, question)
	}

	if len(createdQuestions) == 1 && len(questionsToProcess) == 1 {
		sendJSON(w, createdQuestions[0], http.StatusCreated)
	} else {
		sendJSON(w, createdQuestions, http.StatusCreated)
	}
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
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

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	queryBuilder := applyFilters(db.Model(&models.Question{}), query)
	queryBuilder = queryBuilder.Order(orderBy)

	total, err := countQuestions(queryBuilder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting questions: %v", err), http.StatusInternalServerError)
		return
	}

	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	queryBuilder = queryBuilder.Scopes(models.ApplyQuestionIncludes(includes))
	offset := (page - 1) * limit
	questions, err := findQuestions(queryBuilder, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", err), http.StatusInternalServerError)
		return
	}

	var questionsResp []models.QuestionResponse
	for _, q := range questions {
		questionsResp = append(questionsResp, q.ToResponse())
	}

	sendPaginatedResponse(w, questionsResp, total, limit, page)
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	question, err := findQuestionByID(db, id, includes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question.ToResponse()); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetQuestionsByArray(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	var req IDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, "Invalid JSON body or empty IDs array", http.StatusBadRequest)
		return
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	questions, err := findQuestionsByIDs(db, req.IDs, includes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", err), http.StatusInternalServerError)
		return
	}

	var questionsResp []models.QuestionResponse
	for _, q := range questions {
		questionsResp = append(questionsResp, q.ToResponse())
	}

	sendJSON(w, questionsResp, http.StatusOK)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	log.Printf("Attempting to delete question ID: %s", id)
	rows, err := deleteQuestionByID(db, id)
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

func applyFilters(query *gorm.DB, params map[string][]string) *gorm.DB {
	for param, handler := range filters.QuestionFilters {
		if values, exists := params[param]; exists && len(values) > 0 && values[0] != "" {
			query = handler(values[0], query)
		}
	}
	return query
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
