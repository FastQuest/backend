package questionoption

import (
	"errors"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateQuestionOptions(options []QuestionOption) (int64, error) {
	if r.db == nil {
		return 0, errors.New("database connection not established")
	}
	result := r.db.Create(&options)
	return result.RowsAffected, result.Error
}

func (r *gormRepository) GetQuestionOptionsByQuestionID(questionID string) ([]QuestionOption, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	var options []QuestionOption
	if err := r.db.Where("id_question = ?", questionID).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *gormRepository) GetQuestionOptionsByIDs(ids []uint) ([]QuestionOption, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	var options []QuestionOption
	if err := r.db.Where("id IN ?", ids).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}
