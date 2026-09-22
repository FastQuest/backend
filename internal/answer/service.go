package answer

import (
	"fmt"
)

type service struct {
	repository Repository
}

// NewService builds the default implementation of Service.
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) SendAnswers(requests []CreateAnswerRequest) error {
	for i, req := range requests {
		if req.QuestionOptionID == 0 {
			return fmt.Errorf("questionOptionID at index %d cannot be zero", i)
		}

		if req.QuestionID == 0 {
			return fmt.Errorf("questionID at index %d cannot be zero", i)
		}

		if req.SubmissionID == 0 {
			return fmt.Errorf("submissionID at index %d cannot be zero", i)
		}
	}

	answers := make([]Answer, len(requests))
	for i, req := range requests {
		answers[i] = Answer{
			QuestionOptionID: req.QuestionOptionID,
			QuestionID:       req.QuestionID,
			SubmissionID:     req.SubmissionID,
			IsCorrect:        req.IsCorrect,
		}
	}

	if _, err := s.repository.CreateAnswers(answers); err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}

	return nil
}

func (s *service) GetSubjectPerformance(userID int) ([]SubjectPerformanceResponse, error) {
	performances, err := s.repository.GetUserPerfomace(userID)
	if err != nil {
		return nil, err
	}

	var responses []SubjectPerformanceResponse
	for _, performance := range performances {
		responses = append(responses, toSubjectPerformanceResponse(performance))
	}

	return responses, nil
}

func (s *service) GetOverallPerformance(userID int) (OverallPerformanceResponse, error) {
	performance, err := s.repository.GetUserGeralPerfomace(userID)
	if err != nil {
		return OverallPerformanceResponse{}, err
	}

	return OverallPerformanceResponse{
		TotalAnswers:      performance.TotalAnswers,
		TotalCorrect:      performance.TotalCorrect,
		PercentualCorrect: performance.PercentualCorrect,
	}, nil
}

func toSubjectPerformanceResponse(performance SubjectPerformance) SubjectPerformanceResponse {
	return SubjectPerformanceResponse{
		Subject: SubjectResponse{
			ID:   performance.SubjectID,
			Name: performance.SubjectName,
		},
		TotalAnswers:      performance.TotalAnswers,
		TotalCorrect:      performance.TotalCorrect,
		PercentualCorrect: performance.PercentualCorrect,
	}
}
