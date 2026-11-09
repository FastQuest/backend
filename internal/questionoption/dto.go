package questionoption

type QuestionOptionsBody struct {
	QuestionOptionIDs []uint `json:"question_option_ids"`
}

type CreateRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type Response struct {
	ID         uint   `json:"id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionID uint   `json:"question_id"`
}

type CreateResult struct {
	Count int64
	IDs   []uint
}
