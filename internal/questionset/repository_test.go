package questionset

import (
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGormRepositoryCreatesAndQueriesQuestionSet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &QuestionSet{}, &QuestionSetQuestion{}); err != nil {
		t.Fatalf("migrate question set models: %v", err)
	}

	user := User{Name: "Learner", Email: "learner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	repository := NewRepository(db)
	questionSet := QuestionSet{
		Name:   "Practice",
		Type:   "list",
		UserID: int(user.ID),
	}
	if err := repository.CreateQuestionSet(&questionSet, []int{20, 10}); err != nil {
		t.Fatalf("create question set: %v", err)
	}
	if err := repository.AddQuestionsToSet(questionSet.ID, []uint{7}); err != nil {
		t.Fatalf("append question to set: %v", err)
	}

	ids, err := repository.FindQuestionSetQuestionIDs("1")
	if err != nil {
		t.Fatalf("find linked question IDs: %v", err)
	}
	if !reflect.DeepEqual(ids, []int{20, 10, 7}) {
		t.Fatalf("expected question IDs in position order, got %v", ids)
	}

	isPrivate := false
	filter := ListFilter{
		UserID:    intPointer(int(user.ID)),
		IsPrivate: &isPrivate,
		Statement: "pract",
		OrderBy:   "created_at desc",
	}
	total, err := repository.CountQuestionSets(filter)
	if err != nil {
		t.Fatalf("count question sets: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected one matching set, got %d", total)
	}

	sets, err := repository.FindQuestionSets(filter, []string{"user"}, 0, 10)
	if err != nil {
		t.Fatalf("find question sets: %v", err)
	}
	if len(sets) != 1 || sets[0].User == nil || sets[0].User.Name != "Learner" {
		t.Fatalf("expected matching set with included user, got %#v", sets)
	}
}

func intPointer(value int) *int {
	return &value
}
