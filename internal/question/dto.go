package question

import "flashquest/pkg/models"

type SafeUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type QuestionInput struct {
	Statement            string           `json:"statement"`
	SubjectID            int              `json:"subject_id"`
	UserID               int              `json:"user_id"`
	SourceExamInstanceID *uint            `json:"source_exam_instance_id"`
	Answers              *[]models.Answer `json:"answers"`
}

type IDsRequest struct {
	IDs []uint `json:"ids"`
}
