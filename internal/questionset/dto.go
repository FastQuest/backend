package questionset

import (
	"time"

	"flashquest/internal/question"
)

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
	Questions   []question.CreateRequest `json:"questions"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	User        *UserResponse       `json:"user,omitempty"`
	Questions   []question.Response `json:"questions,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	IsPrivate   bool                `json:"is_private"`
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
