package question

import "time"

type QuestionOptionRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateRequest struct {
	Statement            string                   `json:"statement"`
	SubjectID            int                      `json:"subject_id"`
	UserID               int                      `json:"user_id"`
	SourceExamInstanceID *uint                    `json:"source_exam_instance_id"`
	QuestionOptions      *[]QuestionOptionRequest `json:"question_options"`
}

// CreateResponse mirrors the serialization POST /questions has always
// answered with: the persisted question under its Go field names, with the
// relations never loaded on creation.
type CreateResponse struct {
	ID                   uint                      `json:"ID"`
	CreatedAt            time.Time                 `json:"CreatedAt"`
	UpdatedAt            time.Time                 `json:"UpdatedAt"`
	Statement            string                    `json:"Statement"`
	SubjectID            int                       `json:"SubjectID"`
	Subject              *SubjectResponse          `json:"Subject"`
	UserID               int                       `json:"UserID"`
	User                 *UserResponse             `json:"User"`
	SourceExamInstanceID *uint                     `json:"SourceExamInstanceID"`
	SourceExamInstance   *SourceResponse           `json:"SourceExamInstance"`
	QuestionOptions      *[]QuestionOptionResponse `json:"QuestionOptions"`
}

type SubjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SourceResponse struct {
	ID       uint                   `json:"id"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type QuestionOptionResponse struct {
	ID         uint   `json:"id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionID uint   `json:"question_id"`
}

type Response struct {
	ID              uint                      `json:"id"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	Statement       string                    `json:"statement"`
	Subject         *SubjectResponse          `json:"subject,omitempty"`
	User            *UserResponse             `json:"user,omitempty"`
	Source          *SourceResponse           `json:"source,omitempty"`
	QuestionOptions *[]QuestionOptionResponse `json:"question_options,omitempty"`
}

// QuestionFilterDTO holds the optional filters of the question listing. An
// empty field means the filter was not requested.
type QuestionFilterDTO struct {
	Statement      string
	QuestionOption string
	Subject        string
	Topic          string
	Source         string
	Year           string
	List           string
}

type FilterItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SourceExamInstanceFilterItem struct {
	ID       uint `json:"id"`
	SourceID uint `json:"source_id"`
	Year     uint `json:"year"`
	Edition  uint `json:"edition"`
	Phase    uint `json:"phase"`
}

type FiltersResponse struct {
	Subjects            []FilterItem                   `json:"subjects"`
	Sources             []FilterItem                   `json:"sources"`
	SourceExamInstances []SourceExamInstanceFilterItem `json:"source_exam_instances"`
	Years               []int                          `json:"years"`
}

type IDsRequest struct {
	IDs []uint `json:"ids"`
}
