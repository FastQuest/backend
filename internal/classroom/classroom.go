package classroom

import "errors"

var (
	ErrNotFound             = errors.New("classroom or activity not found")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrDuplicateParticipant = errors.New("user is already a participant")
	ErrQuestionSetEmpty     = errors.New("question set has no questions")
	ErrInvalidQuestionSet   = errors.New("question set not found or has no valid questions")
	ErrQuestionSetForbidden = errors.New("question set is private")
	ErrInactiveClassroom    = errors.New("classroom is not active")
	ErrInactiveActivity     = errors.New("activity is not active")
	ErrNotActiveParticipant = errors.New("user is not an active participant")
	ErrInvalidStatus        = errors.New("invalid classroom status")
)

type Repository interface {
	CreateClassroom(classroom *Classroom) error
	ListClassroomsForUser(userID uint) ([]Classroom, error)
	UpdateClassroomStatus(classroomID uint, status string) error
	AddParticipant(classroomID, userID uint) (*Participant, error)
	ListParticipants(classroomID uint) ([]Participant, error)
	CreateActivity(classroomID, questionSetID, createdBy uint, title string, allowPrivateQuestionSet bool) (*Activity, error)
	ListActivitiesForUser(classroomID, userID uint, allowAdmin bool) ([]Activity, error)
	GetSubmissionContext(activityID, userID uint) (*SubmissionContext, error)
}

type Service interface {
	CreateClassroom(userID uint, role, name string) (ClassroomResponse, error)
	ListClassrooms(userID uint) ([]ClassroomResponse, error)
	UpdateClassroomStatus(actorID uint, role string, classroomID uint, status string) error
	AddParticipant(actorID uint, role string, classroomID, userID uint) (ParticipantResponse, error)
	ListParticipants(actorID uint, role string, classroomID uint) ([]ParticipantResponse, error)
	CreateActivity(actorID uint, role string, classroomID, questionSetID uint, title string) (ActivityResponse, error)
	ListActivities(userID uint, role string, classroomID uint) ([]ActivityDetailResponse, error)
	GetSubmissionContext(activityID, userID uint) (*SubmissionContext, error)
}
