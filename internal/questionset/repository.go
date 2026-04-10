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
