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

// CreateQuestionSet godoc
// @Summary      Creates a new question set
// @Description  Receives question set data and creates it in the database, associating the provided questions.
// @Tags         Question Set
// @Accept       json
// @Produce      json
// @Param        questionSet body NewList true "Question set data"
// @Success      200  {object}  models.QuestionSetResponse
// @Failure      400  {string}  string "Invalid JSON body"
// @Failure      500  {string}  string "Error creating question set"
// @Router       /question-sets [post]
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

// GetQuestionSet godoc
// @Summary      Fetches a question set by ID
// @Description  Returns details of a specific question set.
// @Tags         Question Set
// @Produce      json
// @Param        id       path      int     true  "Question set ID"
// @Param        include  query     string  false "Relationships to include (e.g., user,questions)"
// @Success      200      {object}  models.QuestionSetResponse
// @Failure      404      {string}  string "Error fetching question set"
// @Router       /question-sets/{id} [get]
func GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	questionSet, err := GetQuestionSetByID(id, includes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet.ToResponse())
}

// GetQuestionsFromSet godoc
// @Summary      Lists questions from a set
// @Description  Returns questions associated with a question set. It can return full question objects or just a list of IDs if the query param 'fields=id' is passed.
// @Tags         Question Set
// @Produce      json
// @Param        id      path      int     true  "Question set ID"
// @Param        fields  query     string  false "Fields to return (e.g., id)"
// @Success      200     {array}   models.QuestionResponse
// @Failure      500     {string}  string "Error fetching question set links"
// @Router       /question-sets/{id}/questions [get]
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

// GetLists godoc
// @Summary      Returns registered question sets
// @Description  Fetches question sets with pagination and filtering support (by user, privacy, or search term).
// @Tags         Question Set
// @Produce      json
// @Param        page      query     int     false "Page number" default(1)
// @Param        perPage   query     int     false "Number of items per page" default(10)
// @Param        orderBy   query     string  false "Sorting order" Enums(created_at desc, created_at asc, name asc, name desc) default(created_at desc)
// @Param        userId    query     int     false "Filter by creator user ID"
// @Param        isPrivate query     bool    false "Filter by visibility"
// @Param        statement query     string  false "Search term (name or description)"
// @Param        include   query     string  false "Relationships to include separated by comma (questions and user)"
// @Success      200       {object}  object "Returns a wrapped object containing 'data' and 'pagination'"
// @Failure      400       {string}  string "Invalid param (e.g., malformed userId or isPrivate)"
// @Failure      500       {string}  string "Error fetching lists"
// @Router       /question-sets [get]
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