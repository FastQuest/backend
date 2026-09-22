package answer

import (
	"errors"
)

// ErrNotFound is the domain level translation of a missing record. The
// repository is responsible for converting infrastructure errors (such as
// gorm.ErrRecordNotFound) into this error so the upper layers never depend on
// the persistence technology.
var ErrNotFound = errors.New("answer not found")

// SubjectPerformance is how many answers a user got right on a given subject.
// It is a read projection of the answer domain, free of any gorm or json tag.
type SubjectPerformance struct {
	SubjectID         uint
	SubjectName       string
	TotalAnswers      int
	TotalCorrect      int
	PercentualCorrect float64
}

// OverallPerformance aggregates the same numbers across every subject.
type OverallPerformance struct {
	TotalAnswers      int
	TotalCorrect      int
	PercentualCorrect float64
}

// Repository describes the persistence contract of the answer domain.
// The main entity of this domain is Answer (see models.go).
type Repository interface {
	CreateAnswers(answers []Answer) (int64, error)
	GetUserPerfomace(userID int) ([]SubjectPerformance, error)
	GetUserGeralPerfomace(userID int) (OverallPerformance, error)
}

// Service describes the business contract of the answer domain.
type Service interface {
	SendAnswers(requests []CreateAnswerRequest) error
	GetSubjectPerformance(userID int) ([]SubjectPerformanceResponse, error)
	GetOverallPerformance(userID int) (OverallPerformanceResponse, error)
}
