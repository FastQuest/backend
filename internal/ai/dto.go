package ai

type TestBody struct {
	Text string `json:"text"`
}

type AIAnswerResponse struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type AIQuestionResponse struct {
	Statement string             `json:"statement"`
	Answers   []AIAnswerResponse `json:"answers"`
}

type AIQuestionSetResponse struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Questions   []AIQuestionResponse `json:"questions"`
}
