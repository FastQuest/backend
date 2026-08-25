package answer

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

func (r *Repository) createAnswers(answers *[]models.Answer) (int64, error) {
	result := r.db.Create(answers)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *Repository) GetUserPerfomace(userID int) ([]SubjectPerformance, error) {
	var performances []SubjectPerformance

	err := r.db.Table("submission s").
		Select(`
            sub.*,
            COUNT(a.id) AS total_answers,
            SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) AS total_correct,
            ROUND((SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) * 100.0) / COUNT(a.id), 2) AS percentual_correct
        `).
		Joins("JOIN answer a ON a.submission_id = s.id").
		Joins("JOIN question q ON a.question_id = q.id").
		Joins("JOIN subject sub ON q.subject_id = sub.id").
		Where("s.user_id = ?", userID).
		Where("s.deleted_at IS NULL AND a.deleted_at IS NULL").
		Group("sub.id"). // Agrupamos pela Primary Key do subject
		Order("percentual_correct ASC").
		Scan(&performances).Error

	if err != nil {
		return nil, err
	}

	return performances, nil
}

func (r *Repository) GetUserGeralPerfomace(userID int) (OverallPerformance, error) {
	var performance OverallPerformance

	err := r.db.Table("submission s").
		Select(`
		COUNT(a.id) AS total_answers,
		SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) AS total_correct,
		ROUND((SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) * 100.0) / NULLIF(COUNT(a.id), 0), 2) AS percentual_correct
	`).
		Joins("JOIN answer a ON a.submission_id = s.id").
		Where("s.user_id = ?", userID).
		Where("s.deleted_at IS NULL AND a.deleted_at IS NULL").
		Scan(&performance).Error

	if err != nil {
		return OverallPerformance{}, err
	}

	return performance, nil
}
