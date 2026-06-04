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
