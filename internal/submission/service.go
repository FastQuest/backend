package submission

import (
	"errors"
	"flashquest/internal/answer"
	"flashquest/internal/questionoption"
	"flashquest/pkg/models"
)

func (r *Repository) CreateSubmissionPayload(req CreateSubmissionRequest) (models.Submission, error) {
	optionsIDs := make([]uint, len(req.Answers))

	for i, a := range req.Answers {
		optionsIDs[i] = uint(a.OptionID)
	}

	db := r.db
	if db == nil {
		return models.Submission{}, errors.New("database connection not established")
	}

	options, _ := questionoption.ReadQuestionOptionsByIDArray(db, optionsIDs)

	submission := models.Submission{
		QuestionSetID: req.QuestionSetID,
		UserID:        req.UserID,
		AnswersCount:  len(req.Answers),
		CorrectCount:  0,
	}

	for _, option := range options {
		if option.Is_correct {
			submission.CorrectCount++
		}
	}

	createdSubmission, err := createSubmission(db, &submission)
	if err != nil {
		return models.Submission{}, err
	}

	answers := make([]answer.CreateAnswerRequest, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = answer.CreateAnswerRequest{
			SubmissionID:     createdSubmission.ID,
			QuestionOptionID: a.OptionID,
			QuestionID:       a.QuestionID,
			IsCorrect:        options[i].Is_correct,
		}
	}

	err = answer.NewService(answer.NewRepository(r.db)).SendAnswers(&answers)
	if err != nil {
		return models.Submission{}, err
	}

	return *createdSubmission, nil
}
