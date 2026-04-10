package question

import (
	"errors"
	"flashquest/pkg/models"
	"fmt"
)

func SendQuestions(qs ...*models.Question) error {
	for i, q := range qs {
		if q.Statement == "" {
			return fmt.Errorf("question.Statement at index %d cannot be empty", i)
		}
	}

	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := createQuestions(db, qs); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}
