package submission

import (
	"errors"
	"reflect"
	"testing"

	"flashquest/internal/answer"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
)

type stubRepository struct {
	submission *Submission
	err        error
}

func (r *stubRepository) CreateSubmission(submission *Submission) (*Submission, error) {
	if r.err != nil {
		return nil, r.err
	}
	submission.ID = 17
	r.submission = submission
	return submission, nil
}

func (r *stubRepository) GetSubmissionsByUserID(uint, *uint, int, int) ([]Submission, int64, error) {
	return nil, 0, nil
}

func (r *stubRepository) GetSubmissionByID(uint, string) (*Submission, error) {
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
	options   []questionoption.Response
	err       error
	callCount int
}

func (s *stubQuestionOptionService) CreateQuestionOptions(string, []questionoption.CreateRequest) (questionoption.CreateResult, error) {
	return questionoption.CreateResult{}, nil
}

func (*stubQuestionOptionService) GetQuestionOptionsByQuestionID(string) ([]questionoption.Response, error) {
	return nil, nil
}

func (s *stubQuestionOptionService) GetQuestionOptionsByIDs([]uint) ([]questionoption.Response, error) {
	s.callCount++
	return s.options, s.err
}

type stubQuestionSetService struct {
	questionset.Service
	err      error
	requests []string
}

func (s *stubQuestionSetService) GetQuestionSet(id string, _ []string) (questionset.Response, error) {
	s.requests = append(s.requests, id)
	return questionset.Response{}, s.err
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
	questionSetService := &stubQuestionSetService{}
	service := NewService(repo, answerService, optionService, questionSetService)

	req := CreateSubmissionRequest{QuestionSetID: 4, UserID: 3}
	req.Answers = []AnswerRequest{
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
	if !reflect.DeepEqual(questionSetService.requests, []string{"4"}) {
		t.Fatalf("expected question set validation for ID 4, got %v", questionSetService.requests)
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
	service := NewService(&stubRepository{}, answerService, optionService, &stubQuestionSetService{})

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
	service := NewService(&stubRepository{}, answerService, &stubQuestionOptionService{}, &stubQuestionSetService{})

	req := CreateSubmissionRequest{}
	req.Answers = []AnswerRequest{{QuestionID: 1, OptionID: 2}}

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
	service := NewService(&stubRepository{}, answerService, optionService, &stubQuestionSetService{})

	req := CreateSubmissionRequest{}
	req.Answers = []AnswerRequest{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); !errors.Is(err, answerErr) {
		t.Fatalf("expected answer service error %v, got %v", answerErr, err)
	}
}

func TestCreateSubmissionRejectsMissingQuestionSet(t *testing.T) {
	repository := &stubRepository{}
	answerService := &stubAnswerService{}
	optionService := &stubQuestionOptionService{}
	questionSetService := &stubQuestionSetService{err: questionset.ErrNotFound}
	service := NewService(repository, answerService, optionService, questionSetService)

	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{QuestionSetID: 99})
	if !errors.Is(err, ErrQuestionSetNotFound) {
		t.Fatalf("expected ErrQuestionSetNotFound, got %v", err)
	}
	if repository.submission != nil {
		t.Fatal("expected submission not to be persisted")
	}
	if optionService.callCount != 0 {
		t.Fatal("expected question options not to be queried")
	}
	if answerService.requests != nil {
		t.Fatal("expected answers not to be sent")
	}
}
