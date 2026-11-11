package questionset

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreateQuestionSet godoc
// @Summary      Creates a new question set
// @Description  Receives question set data and creates it in the database, associating the provided questions.
// @Tags         Question Set
// @Accept       json
// @Produce      json
// @Param        questionSet body NewList true "Question set data"
// @Success      200  {object}  Response
// @Failure      400  {string}  string "Invalid JSON body"
// @Failure      500  {string}  string "Error creating question set"
// @Router       /question-sets [post]
func (h *Handler) CreateQuestionSet(w http.ResponseWriter, r *http.Request) {
	var request NewList
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateQuestionSet(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating question set: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, response, http.StatusOK)
}

// GetQuestionSet godoc
// @Summary      Fetches a question set by ID
// @Description  Returns details of a specific question set.
// @Tags         Question Set
// @Produce      json
// @Param        id       path      int     true  "Question set ID"
// @Param        include  query     string  false "Relationships to include (e.g., user,questions)"
// @Success      200      {object}  Response
// @Failure      404      {string}  string "Question set not found"
// @Router       /question-sets/{id} [get]
func (h *Handler) GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	response, err := h.service.GetQuestionSet(id, parseIncludes(r.URL.Query().Get("include")))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "Question set not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question set: %v", err), http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, response, http.StatusOK)
}

// GetQuestionsFromSet godoc
// @Summary      Lists questions from a set
// @Description  Returns questions associated with a question set. It can return full question objects or just a list of IDs if the query param 'fields=id' is passed.
// @Tags         Question Set
// @Produce      json
// @Param        id      path      int     true  "Question set ID"
// @Param        fields  query     string  false "Fields to return (e.g., id)"
// @Success      200     {array}   question.Response
// @Failure      500     {string}  string "Error fetching question set links"
// @Router       /question-sets/{id}/questions [get]
func (h *Handler) GetQuestionsFromSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if r.URL.Query().Get("fields") == "id" {
		questionIDs, err := h.service.GetQuestionIDsFromSet(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching question set links: %v", err), http.StatusInternalServerError)
			return
		}
		writeJSON(w, questionIDs, http.StatusOK)
		return
	}

	questions, err := h.service.GetQuestionsFromSet(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set links: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, questions, http.StatusOK)
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
// @Success      200       {object}  ListResponse
// @Failure      400       {string}  string "Invalid param (e.g., malformed userId or isPrivate)"
// @Failure      500       {string}  string "Error fetching lists"
// @Router       /question-sets [get]
func (h *Handler) GetLists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := parsePositiveInt(query.Get("page"), 1)
	perPage := parsePositiveInt(query.Get("perPage"), 10)
	if perPage > 100 {
		perPage = 10
	}

	filter := ListFilter{
		OrderBy:   query.Get("orderBy"),
		Statement: query.Get("statement"),
	}
	switch filter.OrderBy {
	case "created_at desc", "created_at asc", "name asc", "name desc":
	default:
		filter.OrderBy = "created_at desc"
	}

	if userID := query.Get("userId"); userID != "" {
		value, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		filter.UserID = &value
	}
	if isPrivate := query.Get("isPrivate"); isPrivate != "" {
		value, err := strconv.ParseBool(isPrivate)
		if err != nil {
			http.Error(w, "Invalid isPrivate value", http.StatusBadRequest)
			return
		}
		filter.IsPrivate = &value
	}

	total, err := h.service.CountQuestionSets(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting lists: %v", err), http.StatusInternalServerError)
		return
	}
	questionSets, err := h.service.GetQuestionSets(
		filter,
		parseIncludes(query.Get("include")),
		(page-1)*perPage,
		perPage,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching lists: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, ListResponse{
		Data: questionSets,
		Pagination: PaginationResponse{
			Total:       total,
			PerPage:     perPage,
			CurrentPage: page,
			LastPage:    int(math.Ceil(float64(total) / float64(perPage))),
		},
	}, http.StatusOK)
}

func parseIncludes(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func writeJSON(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
