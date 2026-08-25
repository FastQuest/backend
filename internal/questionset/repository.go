package questionset

import (
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) CreateQuestionSet(questionSet *models.QuestionSet, questionIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(questionSet).Error; err != nil {
			return err
		}
		for index, questionID := range questionIDs {
			link := models.QuestionSetQuestion{
				QuestionSetID: questionSet.ID,
				QuestionID:    questionID,
				Position:      index + 1,
			}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func createQuestionSets(db *gorm.DB, qs []*models.QuestionSet) error {
	return db.Create(qs).Error
}

func createQuestionSetQuestions(db *gorm.DB, qsq []*models.QuestionSetQuestion) error {
	return db.Create(qsq).Error
}

func findQuestionSetByID(db *gorm.DB, id string, includes []string) (*models.QuestionSet, error) {
	var questionSet models.QuestionSet
	if err := db.Scopes(models.ApplyQuestionSetIncludes(includes)).First(&questionSet, id).Error; err != nil {
		return nil, err
	}
	return &questionSet, nil
}
