package classroom

import "time"

type CreateClassroomRequest struct {
	Name string `json:"name"`
}

type AddParticipantRequest struct {
	UserID uint `json:"user_id"`
}

type UpdateClassroomStatusRequest struct {
	Status string `json:"status"`
}

type CreateActivityRequest struct {
	Title         string `json:"title"`
	QuestionSetID uint   `json:"question_set_id"`
}

type ClassroomResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatorID uint      `json:"creator_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ParticipantResponse struct {
	UserID    uint      `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivityResponse struct {
	ID            uint      `json:"id"`
	ClassroomID   uint      `json:"classroom_id"`
	VersionID     uint      `json:"version_id"`
	QuestionSetID uint      `json:"question_set_id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type ActivityDetailResponse struct {
	ID            uint                       `json:"id"`
	ClassroomID   uint                       `json:"classroom_id"`
	VersionID     uint                       `json:"version_id"`
	QuestionSetID uint                       `json:"question_set_id"`
	Title         string                     `json:"title"`
	Status        string                     `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	Questions     []ActivityQuestionResponse `json:"questions"`
}

type ActivityQuestionResponse struct {
	ID        uint                     `json:"id"`
	Position  int                      `json:"position"`
	Statement string                   `json:"statement"`
	Options   []ActivityOptionResponse `json:"options"`
}

type ActivityOptionResponse struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

type SubmissionContext struct {
	ActivityID    uint
	QuestionSetID uint
	Questions     []SnapshotQuestion
}

type SnapshotQuestion struct {
	ID      uint
	Options map[uint]bool
}
