package question

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type stubRepository struct {
	createQuestionFn   func(question *Question) error
	createQuestionsFn  func(questions []*Question) error
	findQuestionsFn    func(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Question, error)
	findQuestionByIDFn func(id string, includes []string) (*Question, error)
	getFilterOptionsFn func() (*FilterOptions, error)

	created []Question
	nextID  uint
}

func (s *stubRepository) CreateQuestion(question *Question) error {
	if s.createQuestionFn != nil {
		if err := s.createQuestionFn(question); err != nil {
			return err
		}
	}
	s.nextID++
	question.ID = s.nextID
	s.created = append(s.created, *question)
	return nil
}

func (s *stubRepository) CreateQuestions(questions []*Question) error {
	if s.createQuestionsFn != nil {
		return s.createQuestionsFn(questions)
	}
	for _, question := range questions {
		s.nextID++
		question.ID = s.nextID
		s.created = append(s.created, *question)
	}
	return nil
}

func (s *stubRepository) CountQuestions(filter QuestionFilterDTO) (int64, error) {
	return 0, nil
}

func (s *stubRepository) FindQuestions(filter QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Question, error) {
	if s.findQuestionsFn == nil {
		return nil, nil
	}
	return s.findQuestionsFn(filter, orderBy, includes, offset, limit)
}

func (s *stubRepository) FindQuestionByID(id string, includes []string) (*Question, error) {
	if s.findQuestionByIDFn == nil {
		return nil, ErrNotFound
	}
	return s.findQuestionByIDFn(id, includes)
}

func (s *stubRepository) FindQuestionsByIDs(ids []uint, includes []string) ([]Question, error) {
	return nil, nil
}

func (s *stubRepository) DeleteQuestionByID(id string) (int64, error) {
	return 0, nil
}

func (s *stubRepository) GetFilterOptions() (*FilterOptions, error) {
	if s.getFilterOptionsFn == nil {
		return &FilterOptions{}, nil
	}
	return s.getFilterOptionsFn()
}

func TestCreateQuestionMapsRequestToEntity(t *testing.T) {
	repo := &stubRepository{}
	svc := NewService(repo)
	sourceID := uint(4)

	created, err := svc.CreateQuestion(CreateRequest{Statement: "S", SubjectID: 2, UserID: 3, SourceExamInstanceID: &sourceID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.created) != 1 {
		t.Fatalf("expected 1 persisted question, got %d", len(repo.created))
	}
	persisted := repo.created[0]
	if persisted.Statement != "S" || persisted.SubjectID != 2 || persisted.UserID != 3 || persisted.SourceExamInstanceID != &sourceID {
		t.Fatalf("unexpected persisted question %+v", persisted)
	}
	if persisted.CreatedAt.IsZero() || persisted.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be set, got %+v", persisted)
	}
	if created.ID != 1 || created.Statement != "S" || created.SourceExamInstanceID != &sourceID {
		t.Fatalf("unexpected response %+v", created)
	}
}

func TestCreateQuestionRequiresStatementAndUser(t *testing.T) {
	for _, req := range []CreateRequest{{UserID: 1}, {Statement: "S"}, {}} {
		repo := &stubRepository{}
		_, err := NewService(repo).CreateQuestion(req)
		if !errors.Is(err, ErrMissingRequiredFields) {
			t.Fatalf("%+v: expected ErrMissingRequiredFields, got %v", req, err)
		}
		if len(repo.created) != 0 {
			t.Fatalf("%+v: expected nothing persisted", req)
		}
	}
}

func TestCreateQuestionsKeepsQuestionsBeforeAnInvalidOne(t *testing.T) {
	repo := &stubRepository{}

	_, err := NewService(repo).CreateQuestions([]CreateRequest{
		{Statement: "A", UserID: 1},
		{Statement: "", UserID: 1},
		{Statement: "C", UserID: 1},
	})
	if !errors.Is(err, ErrBatchMissingRequiredFields) {
		t.Fatalf("expected ErrBatchMissingRequiredFields, got %v", err)
	}
	if len(repo.created) != 1 || repo.created[0].Statement != "A" {
		t.Fatalf("expected only the first question persisted, got %+v", repo.created)
	}
}

func TestCreateQuestionsSingleItemBatchUsesSingleMessage(t *testing.T) {
	_, err := NewService(&stubRepository{}).CreateQuestions([]CreateRequest{{Statement: ""}})
	if !errors.Is(err, ErrMissingRequiredFields) {
		t.Fatalf("expected ErrMissingRequiredFields, got %v", err)
	}
}

func TestCreateQuestionsReturnsPersistenceErrorUnwrapped(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &stubRepository{createQuestionFn: func(*Question) error { return dbErr }}

	_, err := NewService(repo).CreateQuestions([]CreateRequest{{Statement: "A", UserID: 1}})
	if err != dbErr {
		t.Fatalf("expected the raw persistence error, got %v", err)
	}
}

func TestSendQuestionsValidatesEveryStatementBeforeInserting(t *testing.T) {
	repo := &stubRepository{}

	_, err := NewService(repo).SendQuestions([]CreateRequest{{Statement: "A"}, {Statement: ""}})
	if err == nil || err.Error() != "question.Statement at index 1 cannot be empty" {
		t.Fatalf("unexpected error %v", err)
	}
	if len(repo.created) != 0 {
		t.Fatalf("expected nothing persisted, got %+v", repo.created)
	}
}

func TestSendQuestionsReturnsGeneratedIDsInOrder(t *testing.T) {
	repo := &stubRepository{nextID: 10}

	created, err := NewService(repo).SendQuestions([]CreateRequest{{Statement: "A", UserID: 2}, {Statement: "B", UserID: 2}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(created) != 2 || created[0].ID != 11 || created[1].ID != 12 || created[1].Statement != "B" {
		t.Fatalf("unexpected responses %+v", created)
	}
}

func TestSendQuestionsWrapsPersistenceError(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &stubRepository{createQuestionsFn: func([]*Question) error { return dbErr }}

	_, err := NewService(repo).SendQuestions([]CreateRequest{{Statement: "A"}})
	if !errors.Is(err, dbErr) || !strings.HasPrefix(err.Error(), "failed to create question: ") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestGetQuestionMapsEveryRelation(t *testing.T) {
	createdAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &stubRepository{findQuestionByIDFn: func(id string, includes []string) (*Question, error) {
		if id != "5" || !reflect.DeepEqual(includes, []string{"source", "user"}) {
			t.Fatalf("unexpected arguments %q %v", id, includes)
		}
		return &Question{
			ID:        5,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			Statement: "S",
			Subject:   &Subject{ID: 2, Name: "Penal"},
			User:      &User{ID: 3, Name: "Ana", Email: "ana@x"},
			SourceExamInstance: &ExamInstance{
				Edition: 40, Phase: 1, Year: 2023,
				Source: Source{ID: 1, Name: "OAB", Type: "exam"},
			},
			QuestionOptions: &[]QuestionOption{{ID: 9, Text: "a", IsCorrect: true, QuestionID: 5}},
		}, nil
	}}

	resp, err := NewService(repo).GetQuestion("5", []string{"source", "user"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := json.Marshal(resp)
	want := `{"id":5,"created_at":"2025-01-01T00:00:00Z","updated_at":"2025-01-01T00:00:00Z","statement":"S",` +
		`"subject":{"id":2,"name":"Penal"},"user":{"id":3,"name":"Ana","email":"ana@x"},` +
		`"source":{"id":1,"name":"OAB","type":"exam","metadata":{"edition":40,"phase":1,"year":2023}},` +
		`"question_options":[{"id":9,"text":"a","is_correct":true,"question_id":5}]}`
	if string(got) != want {
		t.Fatalf("unexpected json\nwant %s\ngot  %s", want, got)
	}
}

func TestGetQuestionOmitsMissingRelations(t *testing.T) {
	repo := &stubRepository{findQuestionByIDFn: func(string, []string) (*Question, error) {
		return &Question{ID: 1, Statement: "S", QuestionOptions: &[]QuestionOption{}}, nil
	}}

	resp, err := NewService(repo).GetQuestion("1", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Subject != nil || resp.User != nil || resp.Source != nil || resp.QuestionOptions != nil {
		t.Fatalf("expected no relations, got %+v", resp)
	}
}

func TestGetQuestionPropagatesNotFound(t *testing.T) {
	_, err := NewService(&stubRepository{}).GetQuestion("1", nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetQuestionsKeepsEmptyResultNil(t *testing.T) {
	filter := QuestionFilterDTO{Statement: "x", Year: "2024"}
	repo := &stubRepository{findQuestionsFn: func(f QuestionFilterDTO, orderBy string, includes []string, offset, limit int) ([]Question, error) {
		if f != filter || orderBy != "id asc" || offset != 20 || limit != 10 {
			t.Fatalf("unexpected arguments %+v %q %d %d", f, orderBy, offset, limit)
		}
		return []Question{}, nil
	}}

	resp, err := NewService(repo).GetQuestions(filter, "id asc", nil, 20, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil slice, got %#v", resp)
	}
}

func TestGetQuestionFiltersMapsOptions(t *testing.T) {
	repo := &stubRepository{getFilterOptionsFn: func() (*FilterOptions, error) {
		return &FilterOptions{
			Subjects: []FilterOption{{ID: 1, Name: "Civil"}},
			Sources:  []FilterOption{{ID: 2, Name: "OAB"}},
			SourceExamInstances: []SourceExamInstanceFilterOption{
				{ID: 4, SourceID: 2, Year: 2024, Edition: 41, Phase: 1},
			},
			Years: []int{2024, 2023},
		}, nil
	}}

	resp, err := NewService(repo).GetQuestionFilters()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := json.Marshal(resp)
	want := `{"subjects":[{"id":1,"name":"Civil"}],"sources":[{"id":2,"name":"OAB"}],"source_exam_instances":[{"id":4,"source_id":2,"year":2024,"edition":41,"phase":1}],"years":[2024,2023]}`
	if string(got) != want {
		t.Fatalf("unexpected json\nwant %s\ngot  %s", want, got)
	}
}
