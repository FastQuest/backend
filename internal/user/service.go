package user

import (
	"flashquest/internal/user/dto"
)

type service struct {
	repository Repository
}

// NewService builds the default implementation of Service.
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetCurrentUser(userID uint) (dto.Response, error) {
	user, err := s.repository.GetUserByID(userID)
	if err != nil {
		return dto.Response{}, err
	}

	return toResponse(*user), nil
}

func toResponse(user User) dto.Response {
	return dto.Response{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
