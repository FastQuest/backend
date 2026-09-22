package answer

type CreateAnswerRequest struct {
	SubmissionID     uint `json:"submission_id"`
	QuestionID       uint `json:"question_id"`
	QuestionOptionID uint `json:"option_id"`
	IsCorrect        bool `json:"is_correct"`
}

type SubjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SubjectPerformanceResponse struct {
	Subject           SubjectResponse `json:"subject"`
	TotalAnswers      int             `json:"total_answers"`
	TotalCorrect      int             `json:"total_correct"`
	PercentualCorrect float64         `json:"percentual_correct"`
}

type OverallPerformanceResponse struct {
	TotalAnswers      int     `json:"total_answers"`
	TotalCorrect      int     `json:"total_correct"`
	PercentualCorrect float64 `json:"percentual_correct"`
}
