package questionset

import (
	"errors"

	"flashquest/internal/question"
)

var ErrNotFound = errors.New("question set not found")
var ErrDatabaseNotConfigured = errors.New("database connection not established")

type Repository interface {
	CreateQuestionSet(questionSet *QuestionSet, questionIDs []int) error
	AddQuestionsToSet(questionSetID int, questionIDs []uint) error
	FindQuestionSetByID(id string, includes []string) (*QuestionSet, error)
	FindQuestionSetQuestionIDs(id string) ([]int, error)
	CountQuestionSets(filter ListFilter) (int64, error)
	FindQuestionSets(filter ListFilter, includes []string, offset, limit int) ([]QuestionSet, error)
}

type Service interface {
	CreateQuestionSet(request NewList) (Response, error)
	AddQuestionsToSet(questionSetID int, questionIDs []uint) error
	GetQuestionSet(id string, includes []string) (Response, error)
	GetQuestionIDsFromSet(id string) ([]int, error)
	GetQuestionsFromSet(id string) ([]question.Response, error)
	CountQuestionSets(filter ListFilter) (int64, error)
	GetQuestionSets(filter ListFilter, includes []string, offset, limit int) ([]Response, error)
}

type ListFilter struct {
	UserID    *int
	IsPrivate *bool
	Statement string
	OrderBy   string
}

func newQuestionSet(request NewList) QuestionSet {
	return QuestionSet{
		Name:        request.Name,
		Type:        request.Type,
		Description: request.Description,
		UserID:      request.UserID,
		IsPrivate:   request.IsPrivate,
	}
}

func toResponse(questionSet QuestionSet) Response {
	response := Response{
		ID:          questionSet.ID,
		Name:        questionSet.Name,
		Description: questionSet.Description,
		Type:        questionSet.Type,
		CreatedAt:   questionSet.CreatedAt,
		IsPrivate:   questionSet.IsPrivate,
	}

	if questionSet.User != nil {
		response.User = &UserResponse{
			ID:    questionSet.User.ID,
			Name:  questionSet.User.Name,
			Email: questionSet.User.Email,
		}
	}

	return response
}
