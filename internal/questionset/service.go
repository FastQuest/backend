package questionset

import (
	"errors"
	"flashquest/pkg/models"
	"fmt"
)

func (r *Repository) SendQuestionSets(qs ...*models.QuestionSet) error {
	for i, a := range qs {
		if a.Name == "" {
			return fmt.Errorf("questionSet.Name at index %d cannot be empty", i)
		}
		if a.Description == "" {
			return fmt.Errorf("questionSet.Description at index %d cannot be empty", i)
		}
	}

	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := createQuestionSets(db, qs); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func (r *Repository) sendQuestionSetQuestion(qqs ...*models.QuestionSetQuestion) error {
	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := createQuestionSetQuestions(db, qqs); err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func (r *Repository) SendQuestionSetQuestionInternal(qqs ...*models.QuestionSetQuestion) error {
	return r.sendQuestionSetQuestion(qqs...)
}

func (r *Repository) GetQuestionSetByID(id string, includes []string) (*models.QuestionSet, error) {
	db := r.db
	if db == nil {
		return nil, errors.New("database connection not established")
	}

	questionSet, err := findQuestionSetByID(db, id, includes)
	if err != nil {
		return nil, fmt.Errorf("error fetching question set: %w", err)
	}

	return questionSet, nil
}
