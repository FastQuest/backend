package submission

import (
	"errors"

	"flashquest/internal/answer"
)

var (
	ErrQuestionSetNotFound = errors.New("question set not found")
	ErrInvalidSubmission   = errors.New("invalid submission")
)

type Repository interface {
	CreateSubmission(submission *Submission) (*Submission, error)
	CreateSubmissionWithAnswers(submission *Submission, answers []answer.Answer) (*Submission, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Submission, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*Submission, error)
}

type Service interface {
	CreateSubmissionPayload(req CreateSubmissionRequest) (Response, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Response, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*Response, error)
}
