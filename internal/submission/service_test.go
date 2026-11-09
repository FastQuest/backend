package submission

import (
	"errors"
	"reflect"
	"testing"

	"flashquest/internal/answer"
	"flashquest/pkg/models"
)

type stubRepository struct {
	options    []models.QuestionOption
	submission *models.Submission
	err        error
}

func (r *stubRepository) CreateSubmission(submission *models.Submission) (*models.Submission, error) {
	if r.err != nil {
		return nil, r.err
	}
	submission.ID = 17
	r.submission = submission
	return submission, nil
}

func (r *stubRepository) GetQuestionOptionsByIDs([]uint) ([]models.QuestionOption, error) {
	return r.options, r.err
}

func (r *stubRepository) GetSubmissionsByUserID(uint, *uint, int, int) ([]models.Submission, int64, error) {
	return nil, 0, nil
}

func (r *stubRepository) GetSubmissionByID(uint, string) (*models.Submission, error) {
	return nil, nil
}

type stubAnswerService struct {
	requests []answer.CreateAnswerRequest
	err      error
}

func (s *stubAnswerService) SendAnswers(requests []answer.CreateAnswerRequest) error {
	s.requests = requests
	return s.err
}

func (*stubAnswerService) GetSubjectPerformance(int) ([]answer.SubjectPerformanceResponse, error) {
	return nil, nil
}

func (*stubAnswerService) GetOverallPerformance(int) (answer.OverallPerformanceResponse, error) {
	return answer.OverallPerformanceResponse{}, nil
}

func TestCreateSubmissionInjectsAnswerService(t *testing.T) {
	repo := &stubRepository{
		options: []models.QuestionOption{
			{ID: 8, Is_correct: true},
			{ID: 9, Is_correct: false},
		},
	}
	answerService := &stubAnswerService{}
	service := NewService(repo, answerService)

	req := CreateSubmissionRequest{QuestionSetID: 4, UserID: 3}
	req.Answers = []struct {
		QuestionID uint `json:"question_id"`
		OptionID   uint `json:"option_id"`
	}{
		{QuestionID: 11, OptionID: 9},
		{QuestionID: 12, OptionID: 8},
	}

	created, err := service.CreateSubmissionPayload(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID != 17 || created.CorrectCount != 1 || created.AnswersCount != 2 {
		t.Fatalf("unexpected submission: %+v", created)
	}

	want := []answer.CreateAnswerRequest{
		{SubmissionID: 17, QuestionID: 11, QuestionOptionID: 9, IsCorrect: false},
		{SubmissionID: 17, QuestionID: 12, QuestionOptionID: 8, IsCorrect: true},
	}
	if !reflect.DeepEqual(answerService.requests, want) {
		t.Fatalf("expected answers %+v, got %+v", want, answerService.requests)
	}
}

func TestCreateSubmissionReturnsQuestionOptionLookupError(t *testing.T) {
	repoErr := errors.New("options query failed")
	answerService := &stubAnswerService{}
	service := NewService(&stubRepository{err: repoErr}, answerService)

	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected lookup error %v, got %v", repoErr, err)
	}
	if answerService.requests != nil {
		t.Fatal("expected answer service not to be called")
	}
}

func TestCreateSubmissionRejectsMissingQuestionOption(t *testing.T) {
	answerService := &stubAnswerService{}
	service := NewService(&stubRepository{}, answerService)

	req := CreateSubmissionRequest{}
	req.Answers = []struct {
		QuestionID uint `json:"question_id"`
		OptionID   uint `json:"option_id"`
	}{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); err == nil {
		t.Fatal("expected an error for an unknown question option")
	}
	if answerService.requests != nil {
		t.Fatal("expected answer service not to be called")
	}
}

func TestCreateSubmissionReturnsAnswerServiceError(t *testing.T) {
	answerErr := errors.New("answer persistence failed")
	answerService := &stubAnswerService{err: answerErr}
	service := NewService(&stubRepository{
		options: []models.QuestionOption{{ID: 2, Is_correct: true}},
	}, answerService)

	req := CreateSubmissionRequest{}
	req.Answers = []struct {
		QuestionID uint `json:"question_id"`
		OptionID   uint `json:"option_id"`
	}{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); !errors.Is(err, answerErr) {
		t.Fatalf("expected answer service error %v, got %v", answerErr, err)
	}
}
