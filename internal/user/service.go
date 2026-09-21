package user

type service struct {
	repository Repository
}

// NewService builds the default implementation of Service.
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetCurrentUser(userID uint) (Response, error) {
	user, err := s.repository.GetUserByID(userID)
	if err != nil {
		return Response{}, err
	}

	return toResponse(*user), nil
}

func toResponse(user User) Response {
	return Response{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
