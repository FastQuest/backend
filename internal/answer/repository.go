package answer

import (
	"errors"

	"gorm.io/gorm"
)

// subjectPerformanceRow and overallPerformanceRow are the raw shapes the
// aggregate queries scan into. They stay unexported so no gorm tagged struct
// leaks out of this layer.
type subjectPerformanceRow struct {
	Subject           Subject `gorm:"embedded"`
	TotalAnswers      int     `gorm:"column:total_answers"`
	TotalCorrect      int     `gorm:"column:total_correct"`
	PercentualCorrect float64 `gorm:"column:percentual_correct"`
}

type overallPerformanceRow struct {
	TotalAnswers      int     `gorm:"column:total_answers"`
	TotalCorrect      int     `gorm:"column:total_correct"`
	PercentualCorrect float64 `gorm:"column:percentual_correct"`
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository builds the GORM backed implementation of Repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAnswers(answers []Answer) (int64, error) {
	result := r.db.Create(&answers)
	if result.Error != nil {
		return 0, translateError(result.Error)
	}
	return result.RowsAffected, nil
}

func (r *gormRepository) GetUserPerfomace(userID int) ([]SubjectPerformance, error) {
	var rows []subjectPerformanceRow

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
		Scan(&rows).Error

	if err != nil {
		return nil, translateError(err)
	}

	// Kept as a nil slice when there are no rows so the endpoint still
	// serializes to null, exactly as before the refactor.
	var performances []SubjectPerformance
	for _, row := range rows {
		performances = append(performances, SubjectPerformance{
			SubjectID:         row.Subject.ID,
			SubjectName:       row.Subject.Name,
			TotalAnswers:      row.TotalAnswers,
			TotalCorrect:      row.TotalCorrect,
			PercentualCorrect: row.PercentualCorrect,
		})
	}

	return performances, nil
}

func (r *gormRepository) GetUserGeralPerfomace(userID int) (OverallPerformance, error) {
	var row overallPerformanceRow

	err := r.db.Table("submission s").
		Select(`
		COUNT(a.id) AS total_answers,
		SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) AS total_correct,
		ROUND((SUM(CASE WHEN a.is_correct THEN 1 ELSE 0 END) * 100.0) / NULLIF(COUNT(a.id), 0), 2) AS percentual_correct
	`).
		Joins("JOIN answer a ON a.submission_id = s.id").
		Where("s.user_id = ?", userID).
		Where("s.deleted_at IS NULL AND a.deleted_at IS NULL").
		Scan(&row).Error

	if err != nil {
		return OverallPerformance{}, translateError(err)
	}

	return OverallPerformance{
		TotalAnswers:      row.TotalAnswers,
		TotalCorrect:      row.TotalCorrect,
		PercentualCorrect: row.PercentualCorrect,
	}, nil
}

// translateError keeps gorm sentinels from escaping the repository.
func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
