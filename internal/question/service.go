package question

import (
	"fmt"
	"time"
)

type service struct {
	repository Repository
}

// NewService builds the default implementation of Service.
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) CreateQuestion(request CreateRequest) (CreateResponse, error) {
	created, err := s.CreateQuestions([]CreateRequest{request})
	if err != nil {
		return CreateResponse{}, err
	}
	return created[0], nil
}

// CreateQuestions validates and persists each question in order, one at a
// time, so the questions before an invalid or failing one stay created.
func (s *service) CreateQuestions(requests []CreateRequest) ([]CreateResponse, error) {
	var created []CreateResponse

	for _, request := range requests {
		if request.Statement == "" || request.UserID == 0 {
			if len(requests) > 1 {
				return nil, ErrBatchMissingRequiredFields
			}
			return nil, ErrMissingRequiredFields
		}

		question := toQuestion(request)
		if err := s.repository.CreateQuestion(&question); err != nil {
			return nil, err
		}

		created = append(created, toCreateResponse(question))
	}

	return created, nil
}

// SendQuestions is the bulk creation used by the other domains (ai, exam):
// every statement is validated first and then all questions are inserted at
// once. The responses carry the generated IDs in the request order.
func (s *service) SendQuestions(requests []CreateRequest) ([]CreateResponse, error) {
	for i, request := range requests {
		if request.Statement == "" {
			return nil, fmt.Errorf("question.Statement at index %d cannot be empty", i)
		}
	}

	questions := make([]*Question, len(requests))
	for i, request := range requests {
		question := toQuestion(request)
		questions[i] = &question
	}

	if err := s.repository.CreateQuestions(questions); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	created := make([]CreateResponse, len(questions))
	for i, question := range questions {
		created[i] = toCreateResponse(*question)
	}

	return created, nil
}

func (s *service) CountQuestions(filter QuestionFilterDTO) (int64, error) {
	return s.repository.CountQuestions(filter)
}

func (s *service) GetQuestions(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Response, error) {
	questions, err := s.repository.FindQuestions(filter, orderBy, includes, offset, limit)
	if err != nil {
		return nil, err
	}
	return toResponses(questions), nil
}

func (s *service) GetQuestion(id string, includes []string) (Response, error) {
	question, err := s.repository.FindQuestionByID(id, includes)
	if err != nil {
		return Response{}, err
	}
	return toResponse(*question), nil
}

func (s *service) GetQuestionsByIDs(ids []uint, includes []string) ([]Response, error) {
	questions, err := s.repository.FindQuestionsByIDs(ids, includes)
	if err != nil {
		return nil, err
	}
	return toResponses(questions), nil
}

func (s *service) DeleteQuestion(id string) (int64, error) {
	return s.repository.DeleteQuestionByID(id)
}

func (s *service) GetQuestionFilters() (*FiltersResponse, error) {
	options, err := s.repository.GetFilterOptions()
	if err != nil {
		return nil, err
	}

	return &FiltersResponse{
		Subjects:            toFilterItems(options.Subjects),
		Sources:             toFilterItems(options.Sources),
		SourceExamInstances: toSourceExamInstanceFilterItems(options.SourceExamInstances),
		Years:               options.Years,
	}, nil
}

func toSourceExamInstanceFilterItems(instances []SourceExamInstanceFilterOption) []SourceExamInstanceFilterItem {
	items := make([]SourceExamInstanceFilterItem, len(instances))
	for i, instance := range instances {
		items[i] = SourceExamInstanceFilterItem{
			ID:       instance.ID,
			SourceID: instance.SourceID,
			Year:     instance.Year,
			Edition:  instance.Edition,
			Phase:    instance.Phase,
		}
	}
	return items
}

func toQuestion(request CreateRequest) Question {
	return Question{
		Statement:            request.Statement,
		SubjectID:            request.SubjectID,
		UserID:               request.UserID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		SourceExamInstanceID: request.SourceExamInstanceID,
	}
}

func toCreateResponse(question Question) CreateResponse {
	return CreateResponse{
		ID:                   question.ID,
		CreatedAt:            question.CreatedAt,
		UpdatedAt:            question.UpdatedAt,
		Statement:            question.Statement,
		SubjectID:            question.SubjectID,
		UserID:               question.UserID,
		SourceExamInstanceID: question.SourceExamInstanceID,
	}
}

// toResponses keeps a nil slice when there are no questions so the endpoints
// still serialize it to null.
func toResponses(questions []Question) []Response {
	var responses []Response
	for _, question := range questions {
		responses = append(responses, toResponse(question))
	}
	return responses
}

func toResponse(question Question) Response {
	response := Response{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Statement: question.Statement,
	}

	if question.Subject != nil {
		response.Subject = &SubjectResponse{
			ID:   question.Subject.ID,
			Name: question.Subject.Name,
		}
	}

	if question.User != nil {
		response.User = &UserResponse{
			ID:    question.User.ID,
			Name:  question.User.Name,
			Email: question.User.Email,
		}
	}

	if question.QuestionOptions != nil && len(*question.QuestionOptions) > 0 {
		options := make([]QuestionOptionResponse, len(*question.QuestionOptions))
		for i, option := range *question.QuestionOptions {
			options[i] = QuestionOptionResponse{
				ID:         option.ID,
				Text:       option.Text,
				IsCorrect:  option.IsCorrect,
				QuestionID: option.QuestionID,
			}
		}
		response.QuestionOptions = &options
	}

	if question.SourceExamInstance != nil {
		response.Source = &SourceResponse{
			ID:   question.SourceExamInstance.Source.ID,
			Name: question.SourceExamInstance.Source.Name,
			Type: question.SourceExamInstance.Source.Type,
			Metadata: map[string]interface{}{
				"year":    question.SourceExamInstance.Year,
				"edition": question.SourceExamInstance.Edition,
				"phase":   question.SourceExamInstance.Phase,
			},
		}
	}

	return response
}

// toFilterItems preserves a nil input as nil so the JSON output is unchanged.
func toFilterItems(options []FilterOption) []FilterItem {
	if options == nil {
		return nil
	}

	items := make([]FilterItem, len(options))
	for i, option := range options {
		items[i] = FilterItem{ID: option.ID, Name: option.Name}
	}
	return items
}
