package questionset

import (
	"testing"

	"flashquest/internal/question"
)

func TestCreateQuestionSetPersistsRequestAndReturnsResponse(t *testing.T) {
	repository := &questionSetRepositoryStub{}
	service := NewService(repository, &questionServiceStub{})

	response, err := service.CreateQuestionSet(NewList{
		Name:        "Practice",
		Type:        "list",
		Description: "A practice set",
		IsPrivate:   true,
		UserID:      12,
		Questions:   []int{8, 3},
	})
	if err != nil {
		t.Fatalf("CreateQuestionSet returned error: %v", err)
	}
	if response.ID != 5 || response.Name != "Practice" || !response.IsPrivate {
		t.Fatalf("unexpected response: %#v", response)
	}
	if repository.created == nil || repository.created.UserID != 12 || repository.created.Type != "list" {
		t.Fatalf("unexpected persisted question set: %#v", repository.created)
	}
	if len(repository.createdQuestionIDs) != 2 || repository.createdQuestionIDs[0] != 8 || repository.createdQuestionIDs[1] != 3 {
		t.Fatalf("unexpected linked question IDs: %v", repository.createdQuestionIDs)
	}
}

func TestGetQuestionSetUsesQuestionServiceAndPreservesSetOrder(t *testing.T) {
	repository := &questionSetRepositoryStub{
		set: QuestionSet{
			ID:   5,
			Name: "Practice",
			User: &User{ID: 12, Name: "Learner", Email: "learner@example.com"},
		},
		questionIDs: []int{9, 3},
	}
	questionService := &questionServiceStub{
		responses: []question.Response{
			{ID: 3, Statement: "Second"},
			{ID: 9, Statement: "First"},
		},
	}
	service := NewService(repository, questionService)

	response, err := service.GetQuestionSet("5", []string{"user", "questions"})
	if err != nil {
		t.Fatalf("GetQuestionSet returned error: %v", err)
	}
	if response.User == nil || response.User.Name != "Learner" {
		t.Fatalf("expected included user, got %#v", response.User)
	}
	if len(response.Questions) != 2 || response.Questions[0].ID != 9 || response.Questions[1].ID != 3 {
		t.Fatalf("expected questions in set order, got %#v", response.Questions)
	}
	if len(questionService.requestedIDs) != 2 || questionService.requestedIDs[0] != 9 || questionService.requestedIDs[1] != 3 {
		t.Fatalf("question service received unexpected IDs: %v", questionService.requestedIDs)
	}
}

type questionSetRepositoryStub struct {
	created            *QuestionSet
	createdQuestionIDs []int
	set                QuestionSet
	questionIDs        []int
}

func (r *questionSetRepositoryStub) CreateQuestionSet(questionSet *QuestionSet, questionIDs []int) error {
	copy := *questionSet
	copy.ID = 5
	r.created = &copy
	r.createdQuestionIDs = append([]int(nil), questionIDs...)
	questionSet.ID = copy.ID
	return nil
}

func (r *questionSetRepositoryStub) AddQuestionsToSet(int, []uint) error {
	return nil
}

func (r *questionSetRepositoryStub) FindQuestionSetByID(string, []string) (*QuestionSet, error) {
	return &r.set, nil
}

func (r *questionSetRepositoryStub) FindQuestionSetQuestionIDs(string) ([]int, error) {
	return append([]int(nil), r.questionIDs...), nil
}

func (r *questionSetRepositoryStub) CountQuestionSets(ListFilter) (int64, error) {
	return 0, nil
}

func (r *questionSetRepositoryStub) FindQuestionSets(ListFilter, []string, int, int) ([]QuestionSet, error) {
	return nil, nil
}

type questionServiceStub struct {
	question.Service
	responses    []question.Response
	requestedIDs []uint
}

func (s *questionServiceStub) GetQuestionsByIDs(ids []uint, _ []string) ([]question.Response, error) {
	s.requestedIDs = append([]uint(nil), ids...)
	return s.responses, nil
}
