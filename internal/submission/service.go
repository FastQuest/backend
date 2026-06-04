package submission

import (
	"errors"
	"flashquest/internal/platform/database"
	"flashquest/internal/questionoption"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}

func CreateSubmissionPayload(req CreateSubmissionRequest) (models.Submission, error) {
	optionsIDs := make([]uint, len(req.Answers))

	for i, a := range req.Answers {
		optionsIDs[i] = uint(a.OptionID)
	}

	db := getDB()
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

	return *createdSubmission, nil
}
