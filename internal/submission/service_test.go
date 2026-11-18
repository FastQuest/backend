package submission

import (
	"errors"
	"reflect"
	"testing"

	"flashquest/internal/answer"
	"flashquest/internal/classroom"
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
)

type stubRepository struct {
	submission *Submission
	answers    []answer.Answer
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

func (r *stubRepository) CreateSubmissionWithAnswers(submission *Submission, answers []answer.Answer) (*Submission, error) {
	if r.err != nil {
		return nil, r.err
	}
	submission.ID = 17
	r.submission = submission
	for i := range answers {
		answers[i].SubmissionID = submission.ID
	}
	r.answers = answers
	return submission, nil
}

func (r *stubRepository) GetSubmissionsByUserID(uint, *uint, int, int) ([]Submission, int64, error) {
	return nil, 0, nil
}

func (r *stubRepository) GetSubmissionByID(uint, string) (*Submission, error) {
	return nil, nil
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
	err         error
	requests    []string
	questionIDs []int
}

type stubActivityService struct {
	classroom.Service
	context *classroom.SubmissionContext
	err     error
	calls   int
}

func (s *stubActivityService) GetSubmissionContext(uint, uint) (*classroom.SubmissionContext, error) {
	s.calls++
	return s.context, s.err
}

func (s *stubQuestionSetService) GetQuestionSet(id string, _ []string) (questionset.Response, error) {
	s.requests = append(s.requests, id)
	return questionset.Response{}, s.err
}

func (s *stubQuestionSetService) GetQuestionIDsFromSet(string) ([]int, error) {
	return s.questionIDs, nil
}

func TestCreateSubmissionStoresSubmissionAndAnswersAtomically(t *testing.T) {
	repo := &stubRepository{}
	optionService := &stubQuestionOptionService{
		options: []questionoption.Response{
			{ID: 8, QuestionID: 12, IsCorrect: true},
			{ID: 9, QuestionID: 11, IsCorrect: false},
		},
	}
	questionSetService := &stubQuestionSetService{questionIDs: []int{11, 12}}
	service := NewService(repo, optionService, questionSetService, nil)

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
	gotAnswers := make([]answer.CreateAnswerRequest, len(repo.answers))
	for i, item := range repo.answers {
		gotAnswers[i] = answer.CreateAnswerRequest{
			SubmissionID: item.SubmissionID, QuestionID: item.QuestionID,
			QuestionOptionID: item.QuestionOptionID, IsCorrect: item.IsCorrect,
		}
	}
	if !reflect.DeepEqual(gotAnswers, want) {
		t.Fatalf("expected answers %+v, got %+v", want, gotAnswers)
	}
}

func TestCreateSubmissionReturnsQuestionOptionLookupError(t *testing.T) {
	repoErr := errors.New("options query failed")
	optionService := &stubQuestionOptionService{err: repoErr}
	repository := &stubRepository{}
	service := NewService(repository, optionService, &stubQuestionSetService{questionIDs: []int{1}}, nil)

	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected lookup error %v, got %v", repoErr, err)
	}
	if repository.submission != nil {
		t.Fatal("expected submission not to be persisted")
	}
}

func TestCreateSubmissionRejectsMissingQuestionOption(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, &stubQuestionOptionService{}, &stubQuestionSetService{questionIDs: []int{1}}, nil)

	req := CreateSubmissionRequest{}
	req.Answers = []AnswerRequest{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); err == nil {
		t.Fatal("expected an error for an unknown question option")
	}
	if repository.submission != nil {
		t.Fatal("expected submission not to be persisted")
	}
}

func TestCreateSubmissionRejectsDuplicateQuestionAnswers(t *testing.T) {
	repository := &stubRepository{}
	options := &stubQuestionOptionService{options: []questionoption.Response{
		{ID: 2, QuestionID: 1, IsCorrect: true},
	}}
	service := NewService(repository, options, &stubQuestionSetService{questionIDs: []int{1}}, nil)
	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{
		QuestionSetID: 3,
		Answers: []AnswerRequest{
			{QuestionID: 1, OptionID: 2},
			{QuestionID: 1, OptionID: 2},
		},
	})
	if !errors.Is(err, ErrInvalidSubmission) {
		t.Fatalf("expected duplicate question to be rejected, got %v", err)
	}
	if repository.submission != nil {
		t.Fatal("expected invalid submission not to be persisted")
	}
}

func TestCreateSubmissionReturnsPersistenceError(t *testing.T) {
	answerErr := errors.New("submission transaction failed")
	optionService := &stubQuestionOptionService{
		options: []questionoption.Response{{ID: 2, QuestionID: 1, IsCorrect: true}},
	}
	repository := &stubRepository{err: answerErr}
	service := NewService(repository, optionService, &stubQuestionSetService{questionIDs: []int{1}}, nil)

	req := CreateSubmissionRequest{}
	req.Answers = []AnswerRequest{{QuestionID: 1, OptionID: 2}}

	if _, err := service.CreateSubmissionPayload(req); !errors.Is(err, answerErr) {
		t.Fatalf("expected persistence error %v, got %v", answerErr, err)
	}
}

func TestCreateSubmissionRejectsMissingQuestionSet(t *testing.T) {
	repository := &stubRepository{}
	optionService := &stubQuestionOptionService{}
	questionSetService := &stubQuestionSetService{err: questionset.ErrNotFound}
	service := NewService(repository, optionService, questionSetService, nil)

	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{QuestionSetID: 99})
	if !errors.Is(err, ErrQuestionSetNotFound) {
		t.Fatalf("expected ErrQuestionSetNotFound, got %v", err)
	}
}

func TestActivitySubmissionUsesFrozenAnswerKeyAndRequiresParticipantContext(t *testing.T) {
	repository := &stubRepository{}
	activityService := &stubActivityService{
		context: &classroom.SubmissionContext{
			ActivityID: 7, QuestionSetID: 4,
			Questions: []classroom.SnapshotQuestion{{
				ID: 11, Options: map[uint]bool{92: true, 93: false},
			}},
		},
	}
	service := NewService(repository, &stubQuestionOptionService{}, &stubQuestionSetService{}, activityService)
	activityID := uint(7)
	result, err := service.CreateSubmissionPayload(CreateSubmissionRequest{
		ActivityID: &activityID,
		UserID:     20,
		Answers:    []AnswerRequest{{QuestionID: 11, OptionID: 92}},
	})
	if err != nil {
		t.Fatalf("create activity submission: %v", err)
	}
	if result.ActivityID == nil || *result.ActivityID != activityID || result.QuestionSetID != 4 || result.CorrectCount != 1 {
		t.Fatalf("unexpected activity submission: %+v", result)
	}
	if len(repository.answers) != 1 || !repository.answers[0].IsCorrect {
		t.Fatalf("expected answer to be graded against the saved version: %#v", repository.answers)
	}

	repository.submission = nil
	_, err = service.CreateSubmissionPayload(CreateSubmissionRequest{
		ActivityID: &activityID, UserID: 20,
		Answers: []AnswerRequest{{QuestionID: 12, OptionID: 92}},
	})
	if !errors.Is(err, ErrInvalidActivitySubmission) {
		t.Fatalf("expected answer outside activity version to be rejected, got %v", err)
	}
	if repository.submission != nil {
		t.Fatal("expected invalid activity answer not to be persisted")
	}
}

func TestActivitySubmissionRejectsInactiveOrUnauthorizedParticipant(t *testing.T) {
	activityID := uint(7)
	activityService := &stubActivityService{err: classroom.ErrNotActiveParticipant}
	repository := &stubRepository{}
	service := NewService(repository, &stubQuestionOptionService{}, &stubQuestionSetService{}, activityService)

	_, err := service.CreateSubmissionPayload(CreateSubmissionRequest{
		ActivityID: &activityID, UserID: 20,
	})
	if !errors.Is(err, ErrActivityForbidden) {
		t.Fatalf("expected participant authorization failure, got %v", err)
	}
	if repository.submission != nil {
		t.Fatal("expected unauthorized activity submission not to be persisted")
	}
}
