package answer

type CreateAnswerRequest struct {
	SubmissionID     uint `gorm:"not null" json:"submission_id"`
	QuestionID       uint `gorm:"not null" json:"question_id"`
	QuestionOptionID uint `gorm:"not null" json:"option_id"`
	IsCorrect        bool `gorm:"not null" json:"is_correct"`
}

type SubjectPerformance struct {
	SubjectID         int     `gorm:"column:subject_id" json:"subject_id"`
	TotalAnswers      int     `gorm:"column:total_answers" json:"total_answers"`
	TotalCorrect      int     `gorm:"column:total_correct" json:"total_correct"`
	PercentualCorrect float64 `gorm:"column:percentual_correct" json:"percentual_correct"`
}

type SubjectPerformanceResponse struct {
	UserID            int     `json:"user_id"`
	SubjectID         int     `json:"subject_id"`
	TotalAnswers      int     `json:"total_answers"`
	TotalCorrect      int     `json:"total_correct"`
	PercentualCorrect float64 `json:"percentual_correct"`
}
