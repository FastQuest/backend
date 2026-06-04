package questionset

import (
	database "flashquest/internal/platform/database"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
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
