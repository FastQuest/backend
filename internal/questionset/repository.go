package questionset

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateQuestionSet(questionSet *QuestionSet, questionIDs []int) error {
	if r.db == nil {
		return ErrDatabaseNotConfigured
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(questionSet).Error; err != nil {
			return err
		}

		links := make([]QuestionSetQuestion, len(questionIDs))
		for i, questionID := range questionIDs {
			links[i] = QuestionSetQuestion{
				QuestionSetID: questionSet.ID,
				QuestionID:    questionID,
				Position:      i + 1,
			}
		}
		if len(links) > 0 {
			return tx.Create(&links).Error
		}
		return nil
	})
}

func (r *gormRepository) AddQuestionsToSet(questionSetID int, questionIDs []uint) error {
	if len(questionIDs) == 0 {
		return nil
	}
	if r.db == nil {
		return ErrDatabaseNotConfigured
	}

	var lastPosition int
	if err := r.db.Model(&QuestionSetQuestion{}).
		Select("COALESCE(MAX(position), 0)").
		Where("question_set_id = ?", questionSetID).
		Scan(&lastPosition).Error; err != nil {
		return err
	}

	links := make([]QuestionSetQuestion, len(questionIDs))
	for i, questionID := range questionIDs {
		links[i] = QuestionSetQuestion{
			QuestionSetID: questionSetID,
			QuestionID:    int(questionID),
			Position:      lastPosition + i + 1,
		}
	}
	return r.db.Create(&links).Error
}

func (r *gormRepository) FindQuestionSetByID(id string, includes []string) (*QuestionSet, error) {
	if r.db == nil {
		return nil, ErrDatabaseNotConfigured
	}
	var questionSet QuestionSet
	query := r.db
	for _, include := range includes {
		if include == "user" {
			query = query.Preload("User")
		}
	}

	if err := query.Where("id = ?", id).First(&questionSet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &questionSet, nil
}

func (r *gormRepository) FindQuestionSetQuestionIDs(id string) ([]int, error) {
	if r.db == nil {
		return nil, ErrDatabaseNotConfigured
	}
	var links []QuestionSetQuestion
	if err := r.db.Where("question_set_id = ?", id).Order("position ASC").Find(&links).Error; err != nil {
		return nil, err
	}

	questionIDs := make([]int, len(links))
	for i, link := range links {
		questionIDs[i] = link.QuestionID
	}
	return questionIDs, nil
}

func (r *gormRepository) CountQuestionSets(filter ListFilter) (int64, error) {
	if r.db == nil {
		return 0, ErrDatabaseNotConfigured
	}
	var total int64
	if err := applyListFilter(r.db.Model(&QuestionSet{}), filter).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *gormRepository) FindQuestionSets(filter ListFilter, includes []string, offset, limit int) ([]QuestionSet, error) {
	if r.db == nil {
		return nil, ErrDatabaseNotConfigured
	}
	query := applyListFilter(r.db.Model(&QuestionSet{}), filter)
	for _, include := range includes {
		if include == "user" {
			query = query.Preload("User")
		}
	}

	var questionSets []QuestionSet
	orderBy := filter.OrderBy
	switch orderBy {
	case "created_at desc", "created_at asc", "name asc", "name desc":
	default:
		orderBy = "created_at desc"
	}
	if err := query.Order(orderBy).Offset(offset).Limit(limit).Find(&questionSets).Error; err != nil {
		return nil, err
	}
	return questionSets, nil
}

func applyListFilter(query *gorm.DB, filter ListFilter) *gorm.DB {
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.IsPrivate != nil {
		query = query.Where("is_private = ?", *filter.IsPrivate)
	}
	if filter.Statement != "" {
		likeSearch := fmt.Sprintf("%%%s%%", filter.Statement)
		query = query.Where(
			"(LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?))",
			likeSearch, likeSearch,
		)
	}
	return query
}
