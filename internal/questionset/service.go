package questionset

import (
	"errors"
	"flashquest/pkg/models"
	"fmt"
)

func SendQuestionSets(qs ...*models.QuestionSet) error {
	for i, a := range qs {
		if a.Name == "" {
			return fmt.Errorf("questionSet.Name at index %d cannot be empty", i)
		}
		if a.Description == "" {
			return fmt.Errorf("questionSet.Description at index %d cannot be empty", i)
		}
	}

	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := createQuestionSets(db, qs); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func sendQuestionSetQuestion(qqs ...*models.QuestionSetQuestion) error {
	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := createQuestionSetQuestions(db, qqs); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func SendQuestionSetQuestionInternal(qqs ...*models.QuestionSetQuestion) error {
	return sendQuestionSetQuestion(qqs...)
}
