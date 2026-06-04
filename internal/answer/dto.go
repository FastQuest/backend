package answer

type CreateAnswerRequest struct {
	SubmissionID     uint `gorm:"not null" json:"submission_id"`
	QuestionID       uint `gorm:"not null" json:"question_id"`
	QuestionOptionID uint `gorm:"not null" json:"option_id"`
	IsCorrect        bool `gorm:"not null" json:"is_correct"`
}
