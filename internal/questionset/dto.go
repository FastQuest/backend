package questionset

import "flashquest/internal/question"

type NewList struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	UserID      int    `json:"user_id"`
	Questions   []int  `json:"questions"`
}

type NewListWithQuestions struct {
	Name        string                   `json:"name"`
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	IsPrivate   bool                     `json:"is_private"`
	UserID      int                      `json:"user_id"`
	Questions   []question.QuestionInput `json:"questions"`
}
