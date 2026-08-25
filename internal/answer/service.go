package answer

import (
	"flashquest/pkg/models"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) SendAnswers(a *[]CreateAnswerRequest) error {
	for i, ans := range *a {
		if ans.QuestionOptionID == 0 {
			return fmt.Errorf("questionOptionID at index %d cannot be zero", i)
		}

		if ans.QuestionID == 0 {
			return fmt.Errorf("questionID at index %d cannot be zero", i)
		}

		if ans.SubmissionID == 0 {
			return fmt.Errorf("submissionID at index %d cannot be zero", i)
		}
	}

	answers := make([]models.Answer, len(*a))
	for i, ans := range *a {
		answers[i] = models.Answer{
			QuestionOptionID: ans.QuestionOptionID,
			QuestionID:       ans.QuestionID,
			SubmissionID:     ans.SubmissionID,
			IsCorrect:        ans.IsCorrect,
		}
	}

	if _, err := s.repository.createAnswers(&answers); err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}

	return nil
}
