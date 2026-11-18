package submission

import "time"

type AnswerRequest struct {
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

type CreateSubmissionRequest struct {
	QuestionSetID uint            `json:"question_set_id"`
	ActivityID    *uint           `json:"activity_id,omitempty"`
	UserID        uint            `json:"-"`
	Answers       []AnswerRequest `json:"answers"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AnswerResponse struct {
	ID               uint      `json:"id"`
	SubmissionID     uint      `json:"submission_id"`
	QuestionID       uint      `json:"question_id"`
	QuestionOptionID uint      `json:"option_id"`
	IsCorrect        bool      `json:"is_correct"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Response struct {
	ID            uint             `json:"id"`
	UserID        uint             `json:"user_id"`
	QuestionSetID uint             `json:"question_set_id"`
	ActivityID    *uint            `json:"activity_id,omitempty"`
	CorrectCount  int              `json:"correct_count"`
	AnswersCount  int              `json:"answers_count"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	User          *UserResponse    `json:"user,omitempty"`
	Answers       []AnswerResponse `json:"answers,omitempty"`
}

type PaginationResponse struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}

type ListResponse struct {
	Data       []Response         `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
