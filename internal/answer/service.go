package answer

import (
	"errors"
	"flashquest/pkg/models"
	"fmt"
)

func SendAnswers(a *[]models.Answer) error {
	for i, a := range *a {
		if a.Text == "" {
			return fmt.Errorf("answer.Text at index %d cannot be empty", i)
		}
	}

	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if _, err := createAnswers(db, a); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}
