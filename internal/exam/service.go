package exam

import (
	"errors"
	"flashquest/internal/question"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
	"flashquest/internal/source"
	"fmt"
)

func (r *Repository) SendExamInstance(ei ...*source.ExamInstance) error {
	db := r.db
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(ei).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func CreateExamPayload(examRepository *Repository, questionService question.Service, questionOptionService questionoption.Service, questionSetService questionset.Service, newExam NewExam) (questionset.Response, error) {
	exam := newExam.Exam

	errSendE := examRepository.SendExamInstance(&exam)
	if errSendE != nil {
		return questionset.Response{}, errSendE
	}

	questionSet, err := questionSetService.CreateQuestionSet(questionset.NewList{
		Name:        newExam.List.Name,
		Description: newExam.List.Description,
		Type:        "list",
		UserID:      1,
		IsPrivate:   false,
	})
	if err != nil {
		return questionset.Response{}, err
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
		return questionset.Response{}, errSendQ
	}

	for i, q := range newExam.List.Questions {
		var questionoptions []questionoption.CreateRequest
		if q.QuestionOptions != nil {
			questionoptions = make([]questionoption.CreateRequest, 0, len(*q.QuestionOptions))
			for _, qo := range *q.QuestionOptions {
				questionoptions = append(questionoptions, questionoption.CreateRequest{
					Text:      qo.Text,
					IsCorrect: qo.IsCorrect,
				})
			}
		}

		_, errSendQO := questionOptionService.CreateQuestionOptions(fmt.Sprint(questions[i].ID), questionoptions)
		if errSendQO != nil {
			return questionset.Response{}, errSendQO
		}
	}

	questionIDs := make([]uint, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}
	if err := questionSetService.AddQuestionsToSet(questionSet.ID, questionIDs); err != nil {
		return questionset.Response{}, err
	}

	return questionSet, nil
}
