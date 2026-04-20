package question

import (
	database "flashquest/internal/platform/database"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}

func createQuestion(db *gorm.DB, question *models.Question) error {
	return db.Create(question).Error
}

func createQuestions(db *gorm.DB, questions []*models.Question) error {
	return db.Create(questions).Error
}

func countQuestions(qb *gorm.DB) (int64, error) {
	var total int64
	if err := qb.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func findQuestions(qb *gorm.DB, offset, limit int) ([]models.Question, error) {
	var questions []models.Question
	if err := qb.Offset(offset).Limit(limit).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func findQuestionByID(db *gorm.DB, id string, includes []string) (*models.Question, error) {
	var question models.Question
	if err := db.Scopes(models.ApplyQuestionIncludes(includes)).Where("id = ?", id).First(&question).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func findQuestionsByIDs(db *gorm.DB, ids []uint, includes []string) ([]models.Question, error) {
	var questions []models.Question
	if err := db.Scopes(models.ApplyQuestionIncludes(includes)).Where("id IN ?", ids).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func deleteQuestionByID(db *gorm.DB, id string) (int64, error) {
	result := db.Delete(&models.Question{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
