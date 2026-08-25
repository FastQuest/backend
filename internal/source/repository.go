package source

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

func (r *Repository) CreateSource(body models.SourceExamBody) (uint, error) {
	source := models.Source{Name: body.Name, Type: body.Type}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&source).Error; err != nil {
			return err
		}
		examInstance := models.ExamInstance{
			SourceId: source.ID,
			Edition:  body.Edition,
			Phase:    body.Phase,
			Year:     body.Year,
		}
		return tx.Create(&examInstance).Error
	})
	return source.ID, err
}
