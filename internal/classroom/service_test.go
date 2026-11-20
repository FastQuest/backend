package classroom

import (
	"testing"
)

type classroomCreateRepository struct {
	Repository
	classroom *Classroom
}

func (r *classroomCreateRepository) CreateClassroom(classroom *Classroom) error {
	classroom.ID = 1
	r.classroom = classroom
	return nil
}

func TestCreateClassroomAllowsAuthenticatedUser(t *testing.T) {
	repository := &classroomCreateRepository{}
	service := NewService(repository, nil)

	response, err := service.CreateClassroom(6, "Teste")
	if err != nil {
		t.Fatalf("CreateClassroom returned error: %v", err)
	}
	if response.CreatorID != 6 || response.Name != "Teste" || response.Status != StatusActive {
		t.Fatalf("unexpected classroom response: %+v", response)
	}
	if repository.classroom == nil {
		t.Fatal("expected classroom to be persisted")
	}
}
