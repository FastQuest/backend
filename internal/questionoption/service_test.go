package questionoption

import (
	"errors"
	"reflect"
	"testing"

	"flashquest/internal/question"
)

type stubRepository struct {
	created []QuestionOption
	options []QuestionOption
	err     error
}

func (r *stubRepository) CreateQuestionOptions(options []QuestionOption) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	r.created = append([]QuestionOption(nil), options...)
	for i := range options {
		options[i].ID = uint(i + 10)
		r.created[i].ID = options[i].ID
	}
	return int64(len(options)), nil
}

func (r *stubRepository) GetQuestionOptionsByQuestionID(string) ([]QuestionOption, error) {
	return r.options, r.err
}

func (r *stubRepository) GetQuestionOptionsByIDs([]uint) ([]QuestionOption, error) {
	return r.options, r.err
}

type stubQuestionService struct {
	question.Service
	err   error
	calls []string
}

func (s *stubQuestionService) GetQuestion(id string, _ []string) (question.Response, error) {
	s.calls = append(s.calls, id)
	return question.Response{}, s.err
}

func TestCreateQuestionOptionsValidatesQuestionAndMapsRequests(t *testing.T) {
	repo := &stubRepository{}
	questions := &stubQuestionService{}
	service := NewService(repo, questions)

	result, err := service.CreateQuestionOptions("42", []CreateRequest{
		{Text: "A", IsCorrect: false},
		{Text: "B", IsCorrect: true},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(questions.calls, []string{"42"}) {
		t.Fatalf("expected question validation, got calls %v", questions.calls)
	}
	want := []QuestionOption{
		{ID: 10, Text: "A", QuestionID: 42},
		{ID: 11, Text: "B", IsCorrect: true, QuestionID: 42},
	}
	if !reflect.DeepEqual(repo.created, want) {
		t.Fatalf("expected created options %+v, got %+v", want, repo.created)
	}
	if result.Count != 2 || !reflect.DeepEqual(result.IDs, []uint{10, 11}) {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestCreateQuestionOptionsMapsQuestionNotFound(t *testing.T) {
	repo := &stubRepository{}
	questions := &stubQuestionService{err: question.ErrNotFound}

	_, err := NewService(repo, questions).CreateQuestionOptions("404", []CreateRequest{{Text: "A"}})
	if !errors.Is(err, ErrQuestionNotFound) {
		t.Fatalf("expected ErrQuestionNotFound, got %v", err)
	}
	if len(repo.created) != 0 {
		t.Fatalf("expected no options to be persisted, got %+v", repo.created)
	}
}

func TestCreateQuestionOptionsRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name       string
		questionID string
		requests   []CreateRequest
	}{
		{name: "invalid question ID", questionID: "abc"},
		{name: "missing option text", questionID: "1", requests: []CreateRequest{{}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &stubRepository{}
			questions := &stubQuestionService{}
			_, err := NewService(repo, questions).CreateQuestionOptions(tc.questionID, tc.requests)
			if !errors.Is(err, ErrInvalidRequest) {
				t.Fatalf("expected ErrInvalidRequest, got %v", err)
			}
			if len(repo.created) != 0 {
				t.Fatalf("expected no options to be persisted, got %+v", repo.created)
			}
		})
	}
}

func TestGetQuestionOptionsMapsToResponse(t *testing.T) {
	repo := &stubRepository{options: []QuestionOption{
		{ID: 5, Text: "Answer", IsCorrect: true, QuestionID: 9},
	}}

	got, err := NewService(repo, &stubQuestionService{}).GetQuestionOptionsByQuestionID("9")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	want := []Response{{ID: 5, Text: "Answer", IsCorrect: true, QuestionID: 9}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}
