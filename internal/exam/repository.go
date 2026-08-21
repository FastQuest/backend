package exam

import (
	"errors"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetInstanceWithSource(eiID int, ei *models.ExamInstance) error {
	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	result := db.Preload("SourceExamInstance.Source").Where("id = ?", eiID).Find(&ei)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Source Exam Instance not found")
		}
		return errors.New("Error fetching Source Exam Instance")
	}

	return nil
}

