package exam

import (
	"errors"
	"flashquest/internal/question"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
	"flashquest/pkg/sliceutil"
	"fmt"
	"time"
)

func (r *Repository) SendExamInstance(ei ...*models.ExamInstance) error {
	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(ei).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func CreateExamPayload(examRepository *Repository, questionService question.Service, questionOptionService questionoption.Service, questionSetRepository *questionset.Repository, newExam NewExam) (models.QuestionSetResponse, error) {
	exam := newExam.Exam

	errSendE := examRepository.SendExamInstance(&exam)
	if errSendE != nil {
		return models.QuestionSetResponse{}, errSendE
	}

	questionSet := models.QuestionSet{
		Name:        newExam.List.Name,
		Description: newExam.List.Description,
		UserID:      1,
		CreatedAt:   time.Now(),
		IsPrivate:   false,
		Type:        "list",
	}

	errSendQS := questionSetRepository.SendQuestionSets(&questionSet)
	if errSendQS != nil {
		return models.QuestionSetResponse{}, errSendQS
	}

	var questionRequests []question.CreateRequest
	for _, q := range newExam.List.Questions {
		questionRequests = append(questionRequests, question.CreateRequest{
			Statement:            q.Statement,
			SubjectID:            q.SubjectID,
			UserID:               1,
			SourceExamInstanceID: &exam.ID,
		})
	}

	questions, errSendQ := questionService.SendQuestions(questionRequests)
	if errSendQ != nil {
		return models.QuestionSetResponse{}, errSendQ
	}

	for i, q := range newExam.List.Questions {
		questionoptions := make([]questionoption.CreateRequest, 0, len(*q.QuestionOptions))
		for _, qo := range *q.QuestionOptions {
			questionoptions = append(questionoptions, questionoption.CreateRequest{
				Text:      qo.Text,
				IsCorrect: qo.IsCorrect,
			})
		}

		_, errSendQO := questionOptionService.CreateQuestionOptions(fmt.Sprint(questions[i].ID), questionoptions)
		if errSendQO != nil {
			return models.QuestionSetResponse{}, errSendQO
		}
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))
	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := questionSetRepository.SendQuestionSetQuestionInternal(sliceutil.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return models.QuestionSetResponse{}, errSendQSQ
	}

	return questionSet.ToResponse(), nil
}
