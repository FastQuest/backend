package questionset

import (
	"encoding/json"
	"flashquest/pkg/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateQuestionSet(w http.ResponseWriter, r *http.Request) {
	var newList NewList

	err := json.NewDecoder(r.Body).Decode(&newList)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	questionSet := models.QuestionSet{
		Name:        newList.Name,
		Type:        newList.Type,
		Description: newList.Description,
		UserID:      newList.UserID,
		IsPrivate:   newList.IsPrivate,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&questionSet).Error; err != nil {
			return err
		}

		for index, questionID := range newList.Questions {
			link := models.QuestionSetQuestion{
				QuestionSetID: questionSet.ID,
				QuestionID:    questionID,
				Position:      index + 1,
			}

			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating question set: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet.ToResponse())
}

func GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	db := getDB()

	var questionSet models.QuestionSet

	result := db.Scopes(models.ApplyQuestionSetIncludes(includes)).First(&questionSet, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set: %v", result.Error), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet.ToResponse())
}

func GetQuestionsFromSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := getDB()

	query := r.URL.Query()
	returnIDs := query.Get("fields") == "id"

	var links []models.QuestionSetQuestion
	result := db.Where("question_set_id = ?", id).Order("position ASC").Find(&links)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set links: %v", result.Error), http.StatusInternalServerError)
		return
	}

	if len(links) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	w.Header().Set("Content-Type", "application/json")
	if returnIDs {
		json.NewEncoder(w).Encode(questionIDs)
	} else {
		var questions []models.Question
		result = db.Where("id IN ?", questionIDs).Find(&questions)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error), http.StatusInternalServerError)
			return
		}

		var jQuestions []models.QuestionResponse
		for _, q := range questions {
			jQuestions = append(jQuestions, q.ToResponse())
		}

		json.NewEncoder(w).Encode(jQuestions)
	}
}

func GetLists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("perPage"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orderBy := query.Get("orderBy")
	allowedOrders := map[string]bool{
		"created_at desc": true,
		"created_at asc":  true,
		"name asc":        true,
		"name desc":       true,
	}
	if !allowedOrders[orderBy] {
		orderBy = "created_at desc"
	}

	db := getDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	queryBuilder := db.Model(&models.QuestionSet{})

	if userId := query.Get("userId"); userId != "" {
		uid, err := strconv.Atoi(userId)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		queryBuilder = queryBuilder.Where("userId = ?", uid)
	}

	if isPrivate := query.Get("isPrivate"); isPrivate != "" {
		private, err := strconv.ParseBool(isPrivate)
		if err != nil {
			http.Error(w, "Invalid isPrivate value", http.StatusBadRequest)
			return
		}
		queryBuilder = queryBuilder.Where("isPrivate = ?", private)
	}

	if search := query.Get("statement"); search != "" {
		likeSearch := fmt.Sprintf("%%%s%%", search)
		queryBuilder = queryBuilder.Where(
			"(LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?))",
			likeSearch, likeSearch,
		)
	}

	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error counting lists: %v", err), http.StatusInternalServerError)
		return
	}

	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	offset := (page - 1) * limit
	var lists []models.QuestionSet
	result := queryBuilder.Scopes(models.ApplyQuestionSetIncludes(includes)).Order(orderBy).Offset(offset).Limit(limit).Find(&lists)

	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching lists: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var responseLists []models.QuestionSetResponse

	for _, qs := range lists {
		responseLists = append(responseLists, qs.ToResponse())
	}

	response := map[string]interface{}{
		"data": responseLists,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
