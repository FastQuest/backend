package question

import (
	"errors"
	"time"
)

// Question is the GORM persistence model for the question table and the main
// entity of this domain. It carries no json tags on purpose: serialization is
// the job of the DTOs in dto.go.
type Question struct {
	ID                   uint `gorm:"primaryKey"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Statement            string   `gorm:"not null"`
	SubjectID            int      `gorm:"not null"`
	Subject              *Subject `gorm:"foreignKey:SubjectID"`
	UserID               int      `gorm:"not null"`
	User                 *User    `gorm:"foreignKey:UserID"`
	SourceExamInstanceID *uint
	SourceExamInstance   *ExamInstance `gorm:"foreignKey:SourceExamInstanceID"`
	QuestionOptions      *[]QuestionOption
}

func (Question) TableName() string {
	return "question"
}

// ErrNotFound is the domain level translation of a missing question. The
// repository is responsible for converting infrastructure errors (such as
// gorm.ErrRecordNotFound) into this error so the upper layers never depend on
// the persistence technology.
var ErrNotFound = errors.New("question not found")

// Validation errors of the question creation flow. Their messages are the
// ones the API has always answered with.
var (
	ErrMissingRequiredFields      = errors.New("Statement and User ID are required")
	ErrBatchMissingRequiredFields = errors.New("One or more questions are missing Statement or User ID in the batch request")
)

// FilterOption is a selectable value of the question listing filters.
type FilterOption struct {
	ID   uint
	Name string
}

// FilterOptions groups every value the question listing can be filtered by.
type FilterOptions struct {
	Subjects            []FilterOption
	Sources             []FilterOption
	SourceExamInstances []SourceExamInstanceFilterOption
	Years               []int
}

type SourceExamInstanceFilterOption struct {
	ID       uint
	SourceID uint
	Year     uint
	Edition  uint
	Phase    uint
}

// Repository describes the persistence contract of the question domain.
type Repository interface {
	CreateQuestion(question *Question) error
	CreateQuestions(questions []*Question) error
	CountQuestions(filter QuestionFilterDTO) (int64, error)
	FindQuestions(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Question, error)
	FindQuestionByID(id string, includes []string) (*Question, error)
	FindQuestionsByIDs(ids []uint, includes []string) ([]Question, error)
	DeleteQuestionByID(id string) (int64, error)
	GetFilterOptions() (*FilterOptions, error)
}

// Service describes the business contract of the question domain.
type Service interface {
	CreateQuestion(request CreateRequest) (CreateResponse, error)
	CreateQuestions(requests []CreateRequest) ([]CreateResponse, error)
	SendQuestions(requests []CreateRequest) ([]CreateResponse, error)
	CountQuestions(filter QuestionFilterDTO) (int64, error)
	GetQuestions(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Response, error)
	GetQuestion(id string, includes []string) (Response, error)
	GetQuestionsByIDs(ids []uint, includes []string) ([]Response, error)
	DeleteQuestion(id string) (int64, error)
	GetQuestionFilters() (*FiltersResponse, error)
}
