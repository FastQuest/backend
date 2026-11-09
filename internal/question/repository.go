package question

import (
	"errors"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

// NewRepository builds the GORM backed implementation of Repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateQuestion(question *Question) error {
	return translateError(r.db.Create(question).Error)
}

func (r *gormRepository) CreateQuestions(questions []*Question) error {
	return translateError(r.db.Create(questions).Error)
}

func (r *gormRepository) CountQuestions(filter QuestionFilterDTO) (int64, error) {
	var total int64
	if err := applyFilters(r.db.Model(&Question{}), filter).Count(&total).Error; err != nil {
		return 0, translateError(err)
	}
	return total, nil
}

func (r *gormRepository) FindQuestions(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Question, error) {
	var questions []Question
	err := applyFilters(r.db.Model(&Question{}), filter).
		Order(orderBy).
		Scopes(applyIncludes(includes)).
		Offset(offset).
		Limit(limit).
		Find(&questions).Error
	if err != nil {
		return nil, translateError(err)
	}
	return questions, nil
}

func (r *gormRepository) FindQuestionByID(id string, includes []string) (*Question, error) {
	var question Question
	if err := r.db.Scopes(applyIncludes(includes)).Where("id = ?", id).First(&question).Error; err != nil {
		return nil, translateError(err)
	}
	return &question, nil
}

func (r *gormRepository) FindQuestionsByIDs(ids []uint, includes []string) ([]Question, error) {
	var questions []Question
	if err := r.db.Scopes(applyIncludes(includes)).Where("id IN ?", ids).Find(&questions).Error; err != nil {
		return nil, translateError(err)
	}
	return questions, nil
}

func (r *gormRepository) DeleteQuestionByID(id string) (int64, error) {
	result := r.db.Delete(&Question{}, id)
	if result.Error != nil {
		return 0, translateError(result.Error)
	}
	return result.RowsAffected, nil
}

func (r *gormRepository) GetFilterOptions() (*FilterOptions, error) {
	var options FilterOptions

	err := r.db.Table("subject").
		Select("id, name").
		Order("name ASC").
		Scan(&options.Subjects).Error
	if err != nil {
		return nil, translateError(err)
	}

	err = r.db.Table("source").
		Select("id, name").
		Order("name ASC").
		Scan(&options.Sources).Error
	if err != nil {
		return nil, translateError(err)
	}

	err = r.db.Table("source_exam_instance").
		Select("id, source_id, year, edition, phase").
		Order("source_id ASC, year DESC, edition ASC, phase ASC").
		Scan(&options.SourceExamInstances).Error
	if err != nil {
		return nil, translateError(err)
	}

	err = r.db.Table("source_exam_instance").
		Distinct("year").
		Where("year IS NOT NULL").
		Order("year DESC").
		Pluck("year", &options.Years).Error
	if err != nil {
		return nil, translateError(err)
	}

	return &options, nil
}

// translateError keeps gorm sentinels from escaping the repository.
func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
