package answer

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"flashquest/internal/answer/dto"
	"gorm.io/gorm"
)

type stubRepository struct {
	createAnswersFn         func(answers []Answer) (int64, error)
	getUserPerfomaceFn      func(userID int) ([]SubjectPerformance, error)
	getUserGeralPerfomaceFn func(userID int) (OverallPerformance, error)

	createdAnswers []Answer
}

func (s *stubRepository) CreateAnswers(answers []Answer) (int64, error) {
	s.createdAnswers = answers
	if s.createAnswersFn == nil {
		return int64(len(answers)), nil
	}
	return s.createAnswersFn(answers)
}

func (s *stubRepository) GetUserPerfomace(userID int) ([]SubjectPerformance, error) {
	if s.getUserPerfomaceFn == nil {
		return nil, nil
	}
	return s.getUserPerfomaceFn(userID)
}

func (s *stubRepository) GetUserGeralPerfomace(userID int) (OverallPerformance, error) {
	if s.getUserGeralPerfomaceFn == nil {
		return OverallPerformance{}, nil
	}
	return s.getUserGeralPerfomaceFn(userID)
}

func TestSendAnswersMapsRequestsToEntities(t *testing.T) {
	repo := &stubRepository{}
	svc := NewService(repo)

	err := svc.SendAnswers([]dto.CreateAnswerRequest{
		{SubmissionID: 7, QuestionID: 8, QuestionOptionID: 9, IsCorrect: true},
		{SubmissionID: 7, QuestionID: 10, QuestionOptionID: 11, IsCorrect: false},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []Answer{
		{SubmissionID: 7, QuestionID: 8, QuestionOptionID: 9, IsCorrect: true},
		{SubmissionID: 7, QuestionID: 10, QuestionOptionID: 11, IsCorrect: false},
	}
	if !reflect.DeepEqual(repo.createdAnswers, want) {
		t.Fatalf("expected %+v, got %+v", want, repo.createdAnswers)
	}
}

func TestSendAnswersRejectsZeroedIDs(t *testing.T) {
	cases := []struct {
		name    string
		request dto.CreateAnswerRequest
		want    string
	}{
		{"missing option", dto.CreateAnswerRequest{SubmissionID: 1, QuestionID: 1}, "questionOptionID at index 0"},
		{"missing question", dto.CreateAnswerRequest{SubmissionID: 1, QuestionOptionID: 1}, "questionID at index 0"},
		{"missing submission", dto.CreateAnswerRequest{QuestionID: 1, QuestionOptionID: 1}, "submissionID at index 0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &stubRepository{}
			err := NewService(repo).SendAnswers([]dto.CreateAnswerRequest{tc.request})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
			if repo.createdAnswers != nil {
				t.Fatal("expected nothing to be persisted when validation fails")
			}
		})
	}
}

func TestSendAnswersWrapsRepositoryError(t *testing.T) {
	repoErr := errors.New("boom")
	svc := NewService(&stubRepository{
		createAnswersFn: func([]Answer) (int64, error) { return 0, repoErr },
	})

	err := svc.SendAnswers([]dto.CreateAnswerRequest{
		{SubmissionID: 1, QuestionID: 1, QuestionOptionID: 1},
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected the repository error to be wrapped, got %v", err)
	}
}

func TestGetSubjectPerformanceMapsToDTO(t *testing.T) {
	svc := NewService(&stubRepository{
		getUserPerfomaceFn: func(userID int) ([]SubjectPerformance, error) {
			if userID != 42 {
				t.Fatalf("expected user id 42, got %d", userID)
			}
			return []SubjectPerformance{
				{SubjectID: 3, SubjectName: "Direito Civil", TotalAnswers: 10, TotalCorrect: 7, PercentualCorrect: 70},
			}, nil
		},
	})

	got, err := svc.GetSubjectPerformance(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []dto.SubjectPerformanceResponse{
		{
			Subject:           dto.SubjectResponse{ID: 3, Name: "Direito Civil"},
			TotalAnswers:      10,
			TotalCorrect:      7,
			PercentualCorrect: 70,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

// The endpoint used to serialize an empty performance list as null, so the
// service must keep returning a nil slice instead of an empty one.
func TestGetSubjectPerformanceKeepsNilSliceWhenEmpty(t *testing.T) {
	svc := NewService(&stubRepository{
		getUserPerfomaceFn: func(int) ([]SubjectPerformance, error) { return nil, nil },
	})

	got, err := svc.GetSubjectPerformance(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected a nil slice, got %#v", got)
	}
}

func TestGetOverallPerformanceMapsToDTO(t *testing.T) {
	svc := NewService(&stubRepository{
		getUserGeralPerfomaceFn: func(int) (OverallPerformance, error) {
			return OverallPerformance{TotalAnswers: 20, TotalCorrect: 5, PercentualCorrect: 25}, nil
		},
	})

	got, err := svc.GetOverallPerformance(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := dto.OverallPerformanceResponse{TotalAnswers: 20, TotalCorrect: 5, PercentualCorrect: 25}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestTranslateErrorIsolatesGormNotFound(t *testing.T) {
	if err := translateError(gorm.ErrRecordNotFound); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	other := errors.New("connection refused")
	if err := translateError(other); !errors.Is(err, other) {
		t.Fatalf("expected the original error to pass through, got %v", err)
	}
	if translateError(nil) != nil {
		t.Fatal("expected nil to pass through")
	}
}

func TestAnswerTableName(t *testing.T) {
	if table := (Answer{}).TableName(); table != "answer" {
		t.Fatalf("expected answer table, got %s", table)
	}
	if table := (Subject{}).TableName(); table != "subject" {
		t.Fatalf("expected subject table, got %s", table)
	}
}
