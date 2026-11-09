package submission

import (
	"flashquest/internal/answer"
	"flashquest/internal/questionoption"
	"flashquest/pkg/models"
	"fmt"
)

type repository interface {
	CreateSubmission(submission *models.Submission) (*models.Submission, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]models.Submission, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*models.Submission, error)
}

type Service interface {
	CreateSubmissionPayload(req CreateSubmissionRequest) (models.Submission, error)
	GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]models.Submission, int64, error)
	GetSubmissionByID(submissionID uint, includes string) (*models.Submission, error)
}

type service struct {
	repository            repository
	answerService         answer.Service
	questionOptionService questionoption.Service
}

func NewService(repository repository, answerService answer.Service, questionOptionService questionoption.Service) Service {
	return &service{
		repository:            repository,
		answerService:         answerService,
		questionOptionService: questionOptionService,
	}
}

func (s *service) CreateSubmissionPayload(req CreateSubmissionRequest) (models.Submission, error) {
	optionsIDs := make([]uint, len(req.Answers))

	for i, a := range req.Answers {
		optionsIDs[i] = uint(a.OptionID)
	}

	options, err := s.questionOptionService.GetQuestionOptionsByIDs(optionsIDs)
	if err != nil {
		return models.Submission{}, err
	}

	optionsByID := make(map[uint]questionoption.Response, len(options))
	for _, option := range options {
		optionsByID[option.ID] = option
	}

	submission := models.Submission{
		QuestionSetID: req.QuestionSetID,
		UserID:        req.UserID,
		AnswersCount:  len(req.Answers),
		CorrectCount:  0,
	}

	for _, answer := range req.Answers {
		option, exists := optionsByID[uint(answer.OptionID)]
		if !exists {
			return models.Submission{}, fmt.Errorf("question option %d not found", answer.OptionID)
		}
		if option.IsCorrect {
			submission.CorrectCount++
		}
	}

	createdSubmission, err := s.repository.CreateSubmission(&submission)
	if err != nil {
		return models.Submission{}, err
	}

	answers := make([]answer.CreateAnswerRequest, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = answer.CreateAnswerRequest{
			SubmissionID:     createdSubmission.ID,
			QuestionOptionID: a.OptionID,
			QuestionID:       a.QuestionID,
			IsCorrect:        optionsByID[uint(a.OptionID)].IsCorrect,
		}
	}

	err = s.answerService.SendAnswers(answers)
	if err != nil {
		return models.Submission{}, err
	}

	return *createdSubmission, nil
}

func (s *service) GetSubmissionsByUserID(userID uint, questionSetID *uint, page, pageSize int) ([]models.Submission, int64, error) {
	return s.repository.GetSubmissionsByUserID(userID, questionSetID, page, pageSize)
}

func (s *service) GetSubmissionByID(submissionID uint, includes string) (*models.Submission, error) {
	return s.repository.GetSubmissionByID(submissionID, includes)
}
