package exam

import (
	"context"
	"time"

	"flashquest/internal/question"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
	"flashquest/pkg/sliceutil"
)

type service struct {
	examRepository           Repository
	questionRepository       *question.Repository
	questionOptionRepository *questionoption.Repository
	questionSetRepository    *questionset.Repository
}

func NewService(examRepository Repository, questionRepository *question.Repository, questionOptionRepository *questionoption.Repository, questionSetRepository *questionset.Repository) Service {
	return &service{
		examRepository:           examRepository,
		questionRepository:       questionRepository,
		questionOptionRepository: questionOptionRepository,
		questionSetRepository:    questionSetRepository,
	}
}

func (s *service) CreateExamPayload(ctx context.Context, userID uint, newExam NewExam) (models.QuestionSetResponse, error) {
	if err := ctx.Err(); err != nil {
		return models.QuestionSetResponse{}, err
	}

	exam := newExam.Exam

	errSendE := s.examRepository.CreateExamInstance(ctx, &exam)
	if errSendE != nil {
		return models.QuestionSetResponse{}, errSendE
	}

	questionSet := models.QuestionSet{
		Name:        newExam.List.Name,
		Description: newExam.List.Description,
		UserID:      int(userID),
		CreatedAt:   time.Now(),
		IsPrivate:   false,
		Type:        "list",
	}

	errSendQS := s.questionSetRepository.SendQuestionSets(&questionSet)
	if errSendQS != nil {
		return models.QuestionSetResponse{}, errSendQS
	}

	var questions []models.Question
	for _, q := range newExam.List.Questions {
		questions = append(questions, models.Question{
			Statement:            q.Statement,
			SubjectID:            q.SubjectID,
			UserID:               int(userID),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			SourceExamInstanceID: &exam.ID,
		})
	}

	errSendQ := s.questionRepository.SendQuestions(sliceutil.PtrSlice(questions)...)
	if errSendQ != nil {
		return models.QuestionSetResponse{}, errSendQ
	}

	questionoptions := make([]models.QuestionOption, 0, len(questions)*4)
	for i, q := range newExam.List.Questions {
		for _, qo := range *q.QuestionOptions {
			questionoptions = append(questionoptions, models.QuestionOption{
				Text:       qo.Text,
				Is_correct: qo.Is_correct,
				QuestionID: questions[i].ID,
			})
		}
	}

	errSendQO := s.questionOptionRepository.SendQuestionOptions(&questionoptions)
	if errSendQO != nil {
		return models.QuestionSetResponse{}, errSendQO
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))
	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := s.questionSetRepository.SendQuestionSetQuestionInternal(sliceutil.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return models.QuestionSetResponse{}, errSendQSQ
	}

	return questionSet.ToResponse(), nil
}
