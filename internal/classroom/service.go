package classroom

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) Service {
	return &service{repository: repository, db: db}
}

func (s *service) CreateClassroom(userID uint, role, name string) (ClassroomResponse, error) {
	if !canCreate(role) {
		return ClassroomResponse{}, ErrForbidden
	}
	name = strings.TrimSpace(name)
	if userID == 0 || name == "" {
		return ClassroomResponse{}, ErrInvalidRequest
	}
	classroom := Classroom{Name: name, CreatorID: userID, Status: StatusActive}
	if err := s.repository.CreateClassroom(&classroom); err != nil {
		return ClassroomResponse{}, err
	}
	return classroomResponse(classroom), nil
}

func (s *service) ListClassrooms(userID uint) ([]ClassroomResponse, error) {
	classrooms, err := s.repository.ListClassroomsForUser(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]ClassroomResponse, len(classrooms))
	for i, item := range classrooms {
		responses[i] = classroomResponse(item)
	}
	return responses, nil
}

func (s *service) UpdateClassroomStatus(actorID uint, role string, classroomID uint, status string) error {
	if status != StatusActive && status != StatusArchived {
		return ErrInvalidStatus
	}
	if err := s.authorizeManager(actorID, role, classroomID); err != nil {
		return err
	}
	return s.repository.UpdateClassroomStatus(classroomID, status)
}

func (s *service) AddParticipant(actorID uint, role string, classroomID, userID uint) (ParticipantResponse, error) {
	if userID == 0 || classroomID == 0 {
		return ParticipantResponse{}, ErrInvalidRequest
	}
	if err := s.authorizeManager(actorID, role, classroomID); err != nil {
		return ParticipantResponse{}, err
	}
	participant, err := s.repository.AddParticipant(classroomID, userID)
	if err != nil {
		return ParticipantResponse{}, err
	}
	return participantResponse(*participant), nil
}

func (s *service) ListParticipants(actorID uint, role string, classroomID uint) ([]ParticipantResponse, error) {
	if err := s.authorizeManager(actorID, role, classroomID); err != nil {
		return nil, err
	}
	participants, err := s.repository.ListParticipants(classroomID)
	if err != nil {
		return nil, err
	}
	responses := make([]ParticipantResponse, len(participants))
	for i, item := range participants {
		responses[i] = participantResponse(item)
	}
	return responses, nil
}

func (s *service) CreateActivity(actorID uint, role string, classroomID, questionSetID uint, title string) (ActivityResponse, error) {
	title = strings.TrimSpace(title)
	if classroomID == 0 || questionSetID == 0 || title == "" {
		return ActivityResponse{}, ErrInvalidRequest
	}
	if err := s.authorizeManager(actorID, role, classroomID); err != nil {
		return ActivityResponse{}, err
	}
	activity, err := s.repository.CreateActivity(classroomID, questionSetID, actorID, title, role == "Admin")
	if err != nil {
		return ActivityResponse{}, err
	}
	return activityResponse(*activity), nil
}

func (s *service) ListActivities(userID uint, role string, classroomID uint) ([]ActivityDetailResponse, error) {
	activities, err := s.repository.ListActivitiesForUser(classroomID, userID, role == "Admin")
	if err != nil {
		return nil, err
	}
	responses := make([]ActivityDetailResponse, len(activities))
	for i, activity := range activities {
		responses[i] = activityDetailResponse(activity)
	}
	return responses, nil
}

func (s *service) GetSubmissionContext(activityID, userID uint) (*SubmissionContext, error) {
	return s.repository.GetSubmissionContext(activityID, userID)
}

func (s *service) authorizeManager(actorID uint, role string, classroomID uint) error {
	if actorID == 0 {
		return ErrForbidden
	}
	if role != "Admin" && role != "Professor" {
		return ErrForbidden
	}
	var classroom Classroom
	err := s.db.Select("id", "creator_id").First(&classroom, classroomID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if role != "Admin" && classroom.CreatorID != actorID {
		return ErrForbidden
	}
	return nil
}

func canCreate(role string) bool { return role == "Professor" || role == "Admin" }

func classroomResponse(item Classroom) ClassroomResponse {
	return ClassroomResponse{ID: item.ID, Name: item.Name, CreatorID: item.CreatorID, Status: item.Status, CreatedAt: item.CreatedAt}
}

func participantResponse(item Participant) ParticipantResponse {
	return ParticipantResponse{UserID: item.UserID, Status: item.Status, CreatedAt: item.CreatedAt}
}

func activityResponse(item Activity) ActivityResponse {
	return ActivityResponse{
		ID: item.ID, ClassroomID: item.ClassroomID, VersionID: item.VersionID,
		QuestionSetID: item.Version.QuestionSetID, Title: item.Title,
		Status: item.Status, CreatedAt: item.CreatedAt,
	}
}

func activityDetailResponse(item Activity) ActivityDetailResponse {
	response := ActivityDetailResponse{
		ID: item.ID, ClassroomID: item.ClassroomID, VersionID: item.VersionID,
		QuestionSetID: item.Version.QuestionSetID, Title: item.Title,
		Status: item.Status, CreatedAt: item.CreatedAt,
		Questions: make([]ActivityQuestionResponse, len(item.Version.Questions)),
	}
	for i, question := range item.Version.Questions {
		questionResponse := ActivityQuestionResponse{
			ID: question.SourceQuestionID, Position: question.Position,
			Statement: question.Statement,
			Options:   make([]ActivityOptionResponse, len(question.Options)),
		}
		for j, option := range question.Options {
			questionResponse.Options[j] = ActivityOptionResponse{ID: option.SourceOptionID, Text: option.Text}
		}
		response.Questions[i] = questionResponse
	}
	return response
}
