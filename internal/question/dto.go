package question

import "flashquest/pkg/models"

type SafeUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type QuestionInput struct {
	Statement            string                    `json:"statement"`
	SubjectID            int                       `json:"subject_id"`
	UserID               int                       `json:"user_id"`
	SourceExamInstanceID *uint                     `json:"source_exam_instance_id"`
	QuestionOptions      *[]models.QuestionOption `json:"question_options"`
}

type IDsRequest struct {
	IDs []uint `json:"ids"`
}
