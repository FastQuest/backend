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

func GetUserPerfomace(db *gorm.DB, userID int) ([]SubjectPerformance, error) {
	var performances []SubjectPerformance

	err := db.Table("submission s").
		Select(`
			q.subject_id,
			COUNT(a.id) AS total_answers,
			SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) AS total_correct,
			ROUND((SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) * 100.0) / COUNT(a.id), 2) AS percentual_correct
		`).
		Joins("JOIN answer a ON a.submission_id = s.id").
		Joins("JOIN question q ON a.question_id = q.id").
		Where("s.user_id = ?", userID).
		Where("s.deleted_at IS NULL AND a.deleted_at IS NULL").
		Group("q.subject_id").
		Order("percentual_correct ASC").
		Scan(&performances).Error

	if err != nil {
		return nil, err
	}

	return performances, nil
}
