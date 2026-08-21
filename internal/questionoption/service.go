package questionoption

import (
	"errors"
	"flashquest/pkg/models"
	"fmt"
)

func (r *Repository) SendQuestionOptions(qo *[]models.QuestionOption) error {
	for i, qo := range *qo {
		if qo.Text == "" {
			return fmt.Errorf("questionOption.Text at index %d cannot be empty", i)
		}
	}

	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	if _, err := createQuestionOptions(db, qo); err != nil {
		return fmt.Errorf("failed to create question option: %w", err)
	}

	return nil
}
