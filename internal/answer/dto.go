package answer

import "flashquest/pkg/models"

type CreateAnswerRequest struct {
	SubmissionID     uint `gorm:"not null" json:"submission_id"`
	QuestionID       uint `gorm:"not null" json:"question_id"`
	QuestionOptionID uint `gorm:"not null" json:"option_id"`
	IsCorrect        bool `gorm:"not null" json:"is_correct"`
}

type SubjectPerformance struct {
	Subject           models.Subject `gorm:"embedded" json:"subject"`
	TotalAnswers      int            `gorm:"column:total_answers" json:"total_answers"`
	TotalCorrect      int            `gorm:"column:total_correct" json:"total_correct"`
	PercentualCorrect float64        `gorm:"column:percentual_correct" json:"percentual_correct"`
}
