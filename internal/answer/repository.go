package answer

import (
	"flashquest/database"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}

func findQuestionByID(db *gorm.DB, questionID string) (*models.Question, error) {
	var question models.Question
	if err := db.First(&question, questionID).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func createAnswers(db *gorm.DB, answers *[]models.Answer) (int64, error) {
	result := db.Create(answers)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func findAnswersByQuestionID(db *gorm.DB, questionID string) ([]models.Answer, error) {
	var answers []models.Answer
	result := db.Where("id_question = ?", questionID).Find(&answers)
	if result.Error != nil {
		return nil, result.Error
	}
	return answers, nil
}

func readAnswersByIDArray(db *gorm.DB, ids []uint) ([]models.Answer, error) {
	var answers []models.Answer
	resultado := db.Where("id IN (?)", ids).Find(&answers)

	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return answers, nil
		}
		return nil, resultado.Error
	}

	return answers, nil
}
