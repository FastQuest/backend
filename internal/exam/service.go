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

func CreateExamPayload(examRepository *Repository, questionRepository *question.Repository, questionOptionRepository *questionoption.Repository, questionSetRepository *questionset.Repository, newExam NewExam) (models.QuestionSetResponse, error) {
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

	errSendQ := questionRepository.SendQuestions(sliceutil.PtrSlice(questions)...)
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

	errSendQO := questionOptionRepository.SendQuestionOptions(&questionoptions)
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

	errSendQSQ := questionSetRepository.SendQuestionSetQuestionInternal(sliceutil.PtrSlice(questionSetQuestion)...)
	if errSendQSQ != nil {
		return models.QuestionSetResponse{}, errSendQSQ
	}

	return questionSet.ToResponse(), nil
}
