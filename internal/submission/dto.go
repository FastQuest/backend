package submission

type CreateSubmissionRequest struct {
	QuestionSetID uint `json:"question_set_id"`
	UserID        uint `json:"user_id"`
	Answers	   []struct {
		QuestionID uint `json:"question_id"`
		OptionID   uint `json:"option_id"`
	} `json:"answers"`
}
