package exam

import (
	"context"
	"flashquest/pkg/models"
)

// Entidade de domínio (não é DTO)
type Exam struct {
	ID        uint
	SourceId  uint
	Edition   uint
	Phase     uint
	Year      uint
}

// O que o Repositório DEVE fazer
type Repository interface {
	GetInstanceWithSource(eiID int, ei *models.ExamInstance) error
	SendExamInstance(ei ...*models.ExamInstance) error
}

// O que o Service DEVE fazer
type Service interface {
	CreateExamPayload(examRepository *Repository, questionRepository interface{}, questionOptionRepository interface{}, questionSetRepository interface{}, newExam NewExam) (models.QuestionSetResponse, error)
}
