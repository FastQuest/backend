package questionoption

import (
	"errors"
	"fmt"
	"strconv"

	"flashquest/internal/question"
)

type service struct {
	repository      Repository
	questionService question.Service
}

func NewService(repository Repository, questionService question.Service) Service {
	return &service{repository: repository, questionService: questionService}
}

func (s *service) CreateQuestionOptions(questionID string, requests []CreateRequest) (CreateResult, error) {
	if questionID == "" {
		return CreateResult{}, fmt.Errorf("%w: question ID is required", ErrInvalidRequest)
	}
	parsedQuestionID, err := strconv.ParseUint(questionID, 10, 0)
	if err != nil || parsedQuestionID == 0 {
		return CreateResult{}, fmt.Errorf("%w: invalid question ID", ErrInvalidRequest)
	}

	_, err = s.questionService.GetQuestion(questionID, nil)
	if err != nil {
		if errors.Is(err, question.ErrNotFound) {
			return CreateResult{}, ErrQuestionNotFound
		}
		return CreateResult{}, fmt.Errorf("failed to validate question: %w", err)
	}

	options := make([]QuestionOption, len(requests))
	for i, request := range requests {
		if request.Text == "" {
			return CreateResult{}, fmt.Errorf("%w: question option text is required (index %d)", ErrInvalidRequest, i)
		}
		options[i] = QuestionOption{
			Text:      request.Text,
			IsCorrect: request.IsCorrect,
		}
	}

	for i := range options {
		options[i].QuestionID = uint(parsedQuestionID)
	}

	count, err := s.repository.CreateQuestionOptions(options)
	if err != nil {
		return CreateResult{}, fmt.Errorf("failed to create question options: %w", err)
	}

	ids := make([]uint, len(options))
	for i, option := range options {
		ids[i] = option.ID
	}
	return CreateResult{Count: count, IDs: ids}, nil
}

func (s *service) GetQuestionOptionsByQuestionID(questionID string) ([]Response, error) {
	options, err := s.repository.GetQuestionOptionsByQuestionID(questionID)
	if err != nil {
		return nil, err
	}
	return toResponses(options), nil
}

func (s *service) GetQuestionOptionsByIDs(ids []uint) ([]Response, error) {
	options, err := s.repository.GetQuestionOptionsByIDs(ids)
	if err != nil {
		return nil, err
	}
	return toResponses(options), nil
}

func toResponses(options []QuestionOption) []Response {
	if options == nil {
		return nil
	}
	responses := make([]Response, len(options))
	for i, option := range options {
		responses[i] = Response{
			ID:         option.ID,
			Text:       option.Text,
			IsCorrect:  option.IsCorrect,
			QuestionID: option.QuestionID,
		}
	}
	return responses
}
