package submission

import (
	"errors"
	"fmt"

	answerdomain "flashquest/internal/answer"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
)

type service struct {
	repository            Repository
	answerService         answerdomain.Service
	questionOptionService questionoption.Service
	questionSetService    questionset.Service
}

func NewService(repository Repository, answerService answerdomain.Service, questionOptionService questionoption.Service, questionSetService questionset.Service) Service {
	return &service{
		repository:            repository,
		answerService:         answerService,
		questionOptionService: questionOptionService,
		questionSetService:    questionSetService,
	}
}

func (s *service) CreateSubmissionPayload(req CreateSubmissionRequest) (Response, error) {
	_, err := s.questionSetService.GetQuestionSet(fmt.Sprint(req.QuestionSetID), nil)
	if err != nil {
		if errors.Is(err, questionset.ErrNotFound) {
			return Response{}, ErrQuestionSetNotFound
		}
		return Response{}, fmt.Errorf("failed to validate question set: %w", err)
	}

	optionsIDs := make([]uint, len(req.Answers))
	for i, answer := range req.Answers {
		optionsIDs[i] = answer.OptionID
	}

	options, err := s.questionOptionService.GetQuestionOptionsByIDs(optionsIDs)
	if err != nil {
		return Response{}, err
	}

	optionsByID := make(map[uint]questionoption.Response, len(options))
	for _, option := range options {
		optionsByID[option.ID] = option
	}

	submission := Submission{
		QuestionSetID: req.QuestionSetID,
		UserID:        req.UserID,
		AnswersCount:  len(req.Answers),
	}
	for _, answer := range req.Answers {
		option, exists := optionsByID[answer.OptionID]
		if !exists {
			return Response{}, fmt.Errorf("question option %d not found", answer.OptionID)
		}
		if option.IsCorrect {
			submission.CorrectCount++
		}
	}

	createdSubmission, err := s.repository.CreateSubmission(&submission)
	if err != nil {
		return Response{}, err
	}

	answers := make([]answerdomain.CreateAnswerRequest, len(req.Answers))
	for i, answer := range req.Answers {
		answers[i] = answerdomain.CreateAnswerRequest{
			SubmissionID:     createdSubmission.ID,
			QuestionOptionID: answer.OptionID,
			QuestionID:       answer.QuestionID,
			IsCorrect:        optionsByID[answer.OptionID].IsCorrect,
		}
	}
	if err := s.answerService.SendAnswers(answers); err != nil {
		return Response{}, err
	}

	return toResponse(*createdSubmission), nil
}

func (s *service) GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]Response, int64, error) {
	submissions, total, err := s.repository.GetSubmissionsByUserID(userID, questionSetID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return toResponses(submissions), total, nil
}

func (s *service) GetSubmissionByID(submissionID uint, includes string) (*Response, error) {
	submission, err := s.repository.GetSubmissionByID(submissionID, includes)
	if err != nil {
		return nil, err
	}

	response := toResponse(*submission)
	return &response, nil
}

func toResponses(submissions []Submission) []Response {
	if submissions == nil {
		return nil
	}
	responses := make([]Response, len(submissions))
	for i, submission := range submissions {
		responses[i] = toResponse(submission)
	}
	return responses
}

func toResponse(submission Submission) Response {
	response := Response{
		ID:            submission.ID,
		UserID:        submission.UserID,
		QuestionSetID: submission.QuestionSetID,
		CorrectCount:  submission.CorrectCount,
		AnswersCount:  submission.AnswersCount,
		CreatedAt:     submission.CreatedAt,
		UpdatedAt:     submission.UpdatedAt,
	}
	if submission.User != nil {
		response.User = &UserResponse{
			ID:    submission.User.ID,
			Name:  submission.User.Name,
			Email: submission.User.Email,
		}
	}
	if submission.Answers != nil {
		response.Answers = make([]AnswerResponse, len(submission.Answers))
		for i, answer := range submission.Answers {
			response.Answers[i] = AnswerResponse{
				ID:               answer.ID,
				SubmissionID:     answer.SubmissionID,
				QuestionID:       answer.QuestionID,
				QuestionOptionID: answer.QuestionOptionID,
				IsCorrect:        answer.IsCorrect,
				CreatedAt:        answer.CreatedAt,
				UpdatedAt:        answer.UpdatedAt,
			}
		}
	}
	return response
}
