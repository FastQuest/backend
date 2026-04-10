package exam

import (
	"errors"
	"flashquest/helpers"
	"flashquest/internal/answer"
	"flashquest/internal/question"
	"flashquest/internal/questionset"
	"flashquest/pkg/models"
	"fmt"
	"time"
)

func SendExamInstance(ei ...*models.ExamInstance) error {
	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(ei).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func createExamPayload(newExam NewExam) (models.QuestionSetResponse, error) {
	exam := newExam.Exam

	errSendE := SendExamInstance(&exam)
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

	errSendQS := questionset.SendQuestionSets(&questionSet)
	if errSendQS != nil {
		return models.QuestionSetResponse{}, errSendQS
	}

	var questions []models.Question
	for _, q := range newExam.List.Questions {
		questions = append(questions, models.Question{
			Statement:            q.Statement,
			SubjectID:            q.SubjectID,
			UserID:               1,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			SourceExamInstanceID: &exam.ID,
		})
	}

	errSendQ := question.SendQuestions(helpers.PtrSlice(questions)...)
	if errSendQ != nil {
		return models.QuestionSetResponse{}, errSendQ
	}

	answers := make([]models.Answer, 0, len(questions)*4)
	for i, q := range newExam.List.Questions {
		for _, a := range *q.Answers {
			answers = append(answers, models.Answer{
				Text:       a.Text,
				Is_correct: a.Is_correct,
				QuestionID: questions[i].ID,
			})
		}
	}

	errSendA := answer.SendAnswers(&answers)
	if errSendA != nil {
		return models.QuestionSetResponse{}, errSendA
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))
	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := questionset.SendQuestionSetQuestionInternal(helpers.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return models.QuestionSetResponse{}, errSendQSQ
	}

	return questionSet.ToResponse(), nil
}
