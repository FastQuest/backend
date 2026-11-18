package submission

import (
	"errors"
	"fmt"

	answerdomain "flashquest/internal/answer"
	"flashquest/internal/classroom"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
)

type service struct {
	repository            Repository
	questionOptionService questionoption.Service
	questionSetService    questionset.Service
	activityService       classroom.Service
}

var (
	ErrActivityNotFound          = errors.New("activity not found")
	ErrActivityForbidden         = errors.New("user cannot submit to this activity")
	ErrInvalidActivitySubmission = errors.New("answers do not match the activity version")
)

func NewService(repository Repository, questionOptionService questionoption.Service, questionSetService questionset.Service, activityService classroom.Service) Service {
	return &service{
		repository:            repository,
		questionOptionService: questionOptionService,
		questionSetService:    questionSetService,
		activityService:       activityService,
	}
}

func (s *service) CreateSubmissionPayload(req CreateSubmissionRequest) (Response, error) {
	var correctByOption map[uint]bool
	if req.ActivityID != nil {
		if s.activityService == nil {
			return Response{}, ErrActivityNotFound
		}
		context, err := s.activityService.GetSubmissionContext(*req.ActivityID, req.UserID)
		if err != nil {
			if errors.Is(err, classroom.ErrNotFound) || errors.Is(err, classroom.ErrNotActiveParticipant) {
				return Response{}, ErrActivityForbidden
			}
			return Response{}, fmt.Errorf("failed to validate activity: %w", err)
		}
		if req.QuestionSetID != 0 && req.QuestionSetID != context.QuestionSetID {
			return Response{}, ErrInvalidActivitySubmission
		}
		req.QuestionSetID = context.QuestionSetID
		correctByOption = make(map[uint]bool)
		questionOptions := make(map[uint]map[uint]bool, len(context.Questions))
		for _, question := range context.Questions {
			questionOptions[question.ID] = question.Options
		}
		seenQuestions := make(map[uint]struct{}, len(req.Answers))
		for _, item := range req.Answers {
			options, questionExists := questionOptions[item.QuestionID]
			correct, optionExists := options[item.OptionID]
			if !questionExists || !optionExists {
				return Response{}, ErrInvalidActivitySubmission
			}
			if _, duplicate := seenQuestions[item.QuestionID]; duplicate {
				return Response{}, ErrInvalidActivitySubmission
			}
			seenQuestions[item.QuestionID] = struct{}{}
			correctByOption[item.OptionID] = correct
		}
	} else {
		_, err := s.questionSetService.GetQuestionSet(fmt.Sprint(req.QuestionSetID), nil)
		if err != nil {
			if errors.Is(err, questionset.ErrNotFound) {
				return Response{}, ErrQuestionSetNotFound
			}
			return Response{}, fmt.Errorf("failed to validate question set: %w", err)
		}
		questionIDs, err := s.questionSetService.GetQuestionIDsFromSet(fmt.Sprint(req.QuestionSetID))
		if err != nil {
			return Response{}, fmt.Errorf("failed to validate questions in question set: %w", err)
		}
		questionsInSet := make(map[uint]struct{}, len(questionIDs))
		for _, questionID := range questionIDs {
			questionsInSet[uint(questionID)] = struct{}{}
		}
		optionsIDs := make([]uint, len(req.Answers))
		for i, item := range req.Answers {
			optionsIDs[i] = item.OptionID
		}
		options, err := s.questionOptionService.GetQuestionOptionsByIDs(optionsIDs)
		if err != nil {
			return Response{}, err
		}
		optionsByID := make(map[uint]questionoption.Response, len(options))
		for _, option := range options {
			optionsByID[option.ID] = option
		}
		correctByOption = make(map[uint]bool, len(options))
		seenQuestions := make(map[uint]struct{}, len(req.Answers))
		for _, item := range req.Answers {
			if _, questionExists := questionsInSet[item.QuestionID]; !questionExists {
				return Response{}, fmt.Errorf("%w: question %d does not belong to question set %d", ErrInvalidSubmission, item.QuestionID, req.QuestionSetID)
			}
			if _, duplicate := seenQuestions[item.QuestionID]; duplicate {
				return Response{}, fmt.Errorf("%w: question %d appears more than once", ErrInvalidSubmission, item.QuestionID)
			}
			seenQuestions[item.QuestionID] = struct{}{}
			option, exists := optionsByID[item.OptionID]
			if !exists {
				return Response{}, fmt.Errorf("%w: question option %d not found", ErrInvalidSubmission, item.OptionID)
			}
			if option.QuestionID != item.QuestionID {
				return Response{}, fmt.Errorf("%w: question option %d does not belong to question %d", ErrInvalidSubmission, item.OptionID, item.QuestionID)
			}
			correctByOption[item.OptionID] = option.IsCorrect
		}
	}

	submission := Submission{
		QuestionSetID: req.QuestionSetID,
		ActivityID:    req.ActivityID,
		UserID:        req.UserID,
		AnswersCount:  len(req.Answers),
	}
	answers := make([]answerdomain.Answer, len(req.Answers))
	for i, item := range req.Answers {
		if correctByOption[item.OptionID] {
			submission.CorrectCount++
		}
		answers[i] = answerdomain.Answer{
			QuestionID: item.QuestionID, QuestionOptionID: item.OptionID,
			IsCorrect: correctByOption[item.OptionID],
		}
	}

	createdSubmission, err := s.repository.CreateSubmissionWithAnswers(&submission, answers)
	if err != nil {
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
		ActivityID:    submission.ActivityID,
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
