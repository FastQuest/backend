package exam

import (
	"flashquest/internal/questionset"
	"flashquest/internal/source"
)

type NewExam struct {
	Exam source.ExamInstance              `json:"exam"`
	List questionset.NewListWithQuestions `json:"list"`
}
