package exam

import (
	"context"
	"errors"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetInstanceWithSource(ctx context.Context, id uint) (*models.ExamInstance, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}

	var ei models.ExamInstance
	result := r.db.WithContext(ctx).Preload("SourceExamInstance.Source").First(&ei, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Source Exam Instance not found")
		}
		return nil, errors.New("Error fetching Source Exam Instance")
	}

	return &ei, nil
}

func (r *gormRepository) CreateExamInstance(ctx context.Context, ei *models.ExamInstance) error {
	if r.db == nil {
		return errors.New("database connection not established")
	}

	if err := r.db.WithContext(ctx).Create(ei).Error; err != nil {
		return err
	}

	return nil
}
