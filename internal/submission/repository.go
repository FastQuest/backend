package submission

import (
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

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
