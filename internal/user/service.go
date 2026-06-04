package user

import (
	"flashquest/internal/platform/database"
	"flashquest/pkg/models"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetCurrentUser(userID uint) (*models.UserResponse, error) {
	db := database.GetDB()
	user, err := s.repository.GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	response := user.ToResponse()
	return &response, nil
}
