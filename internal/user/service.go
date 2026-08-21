package user

import (
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
	user, err := s.repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	response := user.ToResponse()
	return &response, nil
}
