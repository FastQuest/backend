package submission

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSubmission(submission *Submission) (*Submission, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	if err := r.db.Create(submission).Error; err != nil {
		return nil, err
	}
	return submission, nil
}

func (r *gormRepository) GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Submission, int64, error) {
	if r.db == nil {
		return nil, 0, errors.New("database connection not established")
	}

	var submissions []Submission
	var total int64
	query := r.db.Where("user_id = ?", userID)
	if questionSetID != nil {
		query = query.Where("question_set_id = ?", *questionSetID)
	}
	if err := query.Model(&Submission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}
	return submissions, total, nil
}

func (r *gormRepository) GetSubmissionByID(submissionID uint, includes string) (*Submission, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}

	var submission Submission
	query := r.db.Where("id = ?", submissionID)
	for _, include := range strings.Split(includes, ",") {
		switch strings.TrimSpace(include) {
		case "user":
			query = query.Preload("User")
		case "answers":
			query = query.Preload("Answers")
		}
	}
	if err := query.First(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}
