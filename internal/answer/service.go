package answer

import (
	"errors"
	"flashquest/internal/platform/database"
	"flashquest/pkg/models"
	"fmt"
)

func SendAnswers(a *[]CreateAnswerRequest) error {
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

	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if _, err := createAnswers(db, &answers); err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}

	return nil
}
