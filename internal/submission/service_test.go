package submission

import (
	"errors"
	"reflect"
	"testing"

	"flashquest/internal/answer"
	"flashquest/internal/questionoption"
	"flashquest/pkg/models"
)

type stubRepository struct {
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

type stubQuestionOptionService struct {
	options []questionoption.Response
	err     error
}

func (s *stubQuestionOptionService) CreateQuestionOptions(string, []questionoption.CreateRequest) (questionoption.CreateResult, error) {
	return questionoption.CreateResult{}, nil
}

func (*stubQuestionOptionService) GetQuestionOptionsByQuestionID(string) ([]questionoption.Response, error) {
	return nil, nil
}

func (s *stubQuestionOptionService) GetQuestionOptionsByIDs([]uint) ([]questionoption.Response, error) {
	return s.options, s.err
}

func TestCreateSubmissionInjectsAnswerService(t *testing.T) {
	repo := &stubRepository{}
	optionService := &stubQuestionOptionService{
		options: []questionoption.Response{
			{ID: 8, IsCorrect: true},
			{ID: 9, IsCorrect: false},
		},
	}
	answerService := &stubAnswerService{}
	service := NewService(repo, answerService, optionService)

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
	optionService := &stubQuestionOptionService{err: repoErr}
	service := NewService(&stubRepository{}, answerService, optionService)

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
	service := NewService(&stubRepository{}, answerService, &stubQuestionOptionService{})

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
	optionService := &stubQuestionOptionService{
		options: []questionoption.Response{{ID: 2, IsCorrect: true}},
	}
	service := NewService(&stubRepository{}, answerService, optionService)

	req := CreateSubmissionRequest{}
	req.Answers = []struct {
		QuestionID uint `json:"question_id"`
		OptionID   uint `json:"option_id"`
	}{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); !errors.Is(err, answerErr) {
		t.Fatalf("expected answer service error %v, got %v", answerErr, err)
	}
}
