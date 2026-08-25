package questionoption

import (
	"errors"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

var ErrQuestionNotFound = errors.New("question not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func findQuestionByID(db *gorm.DB, questionID string) (*models.Question, error) {
	var question models.Question
	if err := db.First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	return &question, nil
}

func createQuestionOptions(db *gorm.DB, questionoptions *[]models.QuestionOption) (int64, error) {
	result := db.Create(questionoptions)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func findQuestionOptionsByQuestionID(db *gorm.DB, questionID string) ([]models.QuestionOption, error) {
	var questionoptions []models.QuestionOption
	result := db.Where("id_question = ?", questionID).Find(&questionoptions)
	if result.Error != nil {
		return nil, result.Error
	}
	return questionoptions, nil
}

func ReadQuestionOptionsByIDArray(db *gorm.DB, ids []uint) ([]models.QuestionOption, error) {
	var questionoptions []models.QuestionOption
	resultado := db.Where("id IN (?)", ids).Find(&questionoptions)

	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return questionoptions, nil
		}
		return nil, resultado.Error
	}

	return questionoptions, nil
}
