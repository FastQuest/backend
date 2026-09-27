package exam

import (
	"context"
	"flashquest/pkg/models"
)

type Exam struct {
	ID       uint
	SourceId uint
	Edition  uint
	Phase    uint
	Year     uint
}

type Repository interface {
	GetInstanceWithSource(ctx context.Context, id uint) (*models.ExamInstance, error)
	CreateExamInstance(ctx context.Context, ei *models.ExamInstance) error
}

type Service interface {
	CreateExamPayload(ctx context.Context, userID uint, newExam NewExam) (models.QuestionSetResponse, error)
}
