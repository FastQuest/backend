package questionset

import (
	"fmt"
	"strings"

	"flashquest/internal/question"
)

type service struct {
	repository      Repository
	questionService question.Service
}

func NewService(repository Repository, questionService question.Service) Service {
	return &service{
		repository:      repository,
		questionService: questionService,
	}
}

func (s *service) CreateQuestionSet(request NewList) (Response, error) {
	questionSet := newQuestionSet(request)
	if err := s.repository.CreateQuestionSet(&questionSet, request.Questions); err != nil {
		return Response{}, err
	}
	return toResponse(questionSet), nil
}

func (s *service) AddQuestionsToSet(questionSetID int, questionIDs []uint) error {
	return s.repository.AddQuestionsToSet(questionSetID, questionIDs)
}

func (s *service) GetQuestionSet(id string, includes []string) (Response, error) {
	includes = cleanIncludes(includes)
	questionSet, err := s.repository.FindQuestionSetByID(id, includes)
	if err != nil {
		return Response{}, err
	}

	response := toResponse(*questionSet)
	if hasInclude(includes, "questions") {
		questionIDs, err := s.repository.FindQuestionSetQuestionIDs(id)
		if err != nil {
			return Response{}, err
		}
		response.Questions, err = s.getQuestions(questionIDs)
		if err != nil {
			return Response{}, err
		}
	}
	return response, nil
}

func (s *service) GetQuestionIDsFromSet(id string) ([]int, error) {
	return s.repository.FindQuestionSetQuestionIDs(id)
}

func (s *service) GetQuestionsFromSet(id string) ([]question.Response, error) {
	questionIDs, err := s.repository.FindQuestionSetQuestionIDs(id)
	if err != nil {
		return nil, err
	}
	return s.getQuestions(questionIDs)
}

func (s *service) CountQuestionSets(filter ListFilter) (int64, error) {
	return s.repository.CountQuestionSets(filter)
}

func (s *service) GetQuestionSets(filter ListFilter, includes []string, offset, limit int) ([]Response, error) {
	includes = cleanIncludes(includes)
	questionSets, err := s.repository.FindQuestionSets(filter, includes, offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]Response, len(questionSets))
	for i, questionSet := range questionSets {
		responses[i] = toResponse(questionSet)
		if hasInclude(includes, "questions") {
			questionIDs, err := s.repository.FindQuestionSetQuestionIDs(fmt.Sprint(questionSet.ID))
			if err != nil {
				return nil, err
			}
			responses[i].Questions, err = s.getQuestions(questionIDs)
			if err != nil {
				return nil, err
			}
		}
	}
	return responses, nil
}

func (s *service) getQuestions(questionIDs []int) ([]question.Response, error) {
	if len(questionIDs) == 0 {
		return []question.Response{}, nil
	}

	ids := make([]uint, len(questionIDs))
	for i, id := range questionIDs {
		ids[i] = uint(id)
	}
	questions, err := s.questionService.GetQuestionsByIDs(ids, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch questions in question set: %w", err)
	}

	byID := make(map[uint]question.Response, len(questions))
	for _, item := range questions {
		byID[item.ID] = item
	}
	ordered := make([]question.Response, 0, len(questions))
	for _, id := range questionIDs {
		if item, exists := byID[uint(id)]; exists {
			ordered = append(ordered, item)
		}
	}
	return ordered, nil
}

func cleanIncludes(includes []string) []string {
	cleaned := make([]string, 0, len(includes))
	for _, include := range includes {
		if include = strings.TrimSpace(include); include != "" {
			cleaned = append(cleaned, include)
		}
	}
	return cleaned
}

func hasInclude(includes []string, target string) bool {
	for _, include := range includes {
		if include == target {
			return true
		}
	}
	return false
}
