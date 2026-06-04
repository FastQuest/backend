package answer

import (
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func createAnswers(db *gorm.DB, answers *[]models.Answer) (int64, error) {
	result := db.Create(answers)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
