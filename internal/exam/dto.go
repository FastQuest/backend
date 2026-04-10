package exam

import (
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
)

type NewExam struct {
	Exam models.ExamInstance              `json:"exam"`
	List questionset.NewListWithQuestions `json:"list"`
}
