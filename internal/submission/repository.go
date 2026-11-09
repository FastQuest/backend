package submission

import (
	"errors"
	"flashquest/internal/questionoption"
	"flashquest/pkg/models"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSubmission(submission *models.Submission) (*models.Submission, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	return createSubmission(r.db, submission)
}

func (r *Repository) GetQuestionOptionsByIDs(ids []uint) ([]models.QuestionOption, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	return questionoption.ReadQuestionOptionsByIDArray(r.db, ids)
}

func (r *Repository) GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]models.Submission, int64, error) {
	if r.db == nil {
		return nil, 0, errors.New("database connection not established")
	}
	return getSubmissionsByUserID(r.db, userID, questionSetID, page, pageSize)
}

func (r *Repository) GetSubmissionByID(submissionID uint, includes string) (*models.Submission, error) {
	if r.db == nil {
		return nil, errors.New("database connection not established")
	}
	return getSubmissionByID(r.db, submissionID, includes)
}

func createSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
	result := db.Create(submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return submission, nil
}

func getSubmissionsByUserID(db *gorm.DB, userID uint, questionSetID *uint, page, pageSize int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	query := db.Where("user_id = ?", userID)

	if questionSetID != nil {
		query = query.Where("question_set_id = ?", *questionSetID)
	}

	// Count total records
	if err := query.Model(&models.Submission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

func getSubmissionByID(db *gorm.DB, submissionID uint, includes string) (*models.Submission, error) {
	var submission models.Submission

	query := db.Where("id = ?", submissionID)

	// Parse includes parameter
	if includes != "" {
		includeList := strings.Split(includes, ",")
		for _, inc := range includeList {
			inc = strings.TrimSpace(inc)
			switch inc {
			case "user":
				query = query.Preload("User")
			case "answers":
				query = query.Preload("Answers")
			}
		}
	}

	if err := query.First(&submission).Error; err != nil {
		return nil, err
	}

	return &submission, nil
}
