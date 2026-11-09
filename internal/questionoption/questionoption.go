package questionoption

import (
	"errors"
)

var ErrQuestionNotFound = errors.New("question not found")
var ErrInvalidRequest = errors.New("invalid question option request")

type Repository interface {
	CreateQuestionOptions(options []QuestionOption) (int64, error)
	GetQuestionOptionsByQuestionID(questionID string) ([]QuestionOption, error)
	GetQuestionOptionsByIDs(ids []uint) ([]QuestionOption, error)
}

type Service interface {
	CreateQuestionOptions(questionID string, requests []CreateRequest) (CreateResult, error)
	GetQuestionOptionsByQuestionID(questionID string) ([]Response, error)
	GetQuestionOptionsByIDs(ids []uint) ([]Response, error)
}
