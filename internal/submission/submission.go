package submission

import (
	"errors"
)

var ErrQuestionSetNotFound = errors.New("question set not found")

type Repository interface {
	CreateSubmission(submission *Submission) (*Submission, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Submission, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*Submission, error)
}

type Service interface {
	CreateSubmissionPayload(req CreateSubmissionRequest) (Response, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Response, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*Response, error)
}
