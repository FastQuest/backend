package submission

import (
	"testing"

	"flashquest/internal/answer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGormRepositoryPreloadsSubmissionIncludes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&User{}, &Submission{}, &Answer{}); err != nil {
		t.Fatalf("migrate submission models: %v", err)
	}

	user := User{Name: "Learner", Email: "learner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	record := Submission{
		UserID:        user.ID,
		QuestionSetID: 8,
		CorrectCount:  1,
		AnswersCount:  1,
	}
	repository := NewRepository(db)
	created, err := repository.CreateSubmission(&record)
	if err != nil {
		t.Fatalf("create submission: %v", err)
	}
	if err := db.Create(&Answer{
		SubmissionID:     created.ID,
		QuestionID:       12,
		QuestionOptionID: 34,
		IsCorrect:        true,
	}).Error; err != nil {
		t.Fatalf("create submission answer: %v", err)
	}

	got, err := repository.GetSubmissionByID(created.ID, "user, answers")
	if err != nil {
		t.Fatalf("get submission with includes: %v", err)
	}
	if got.User == nil || got.User.Name != "Learner" {
		t.Fatalf("expected included user, got %#v", got.User)
	}
	if len(got.Answers) != 1 || got.Answers[0].QuestionOptionID != 34 {
		t.Fatalf("expected included answer, got %#v", got.Answers)
	}

	submissions, total, err := repository.GetSubmissionsByUserID(user.ID, nil, 1, 10)
	if err != nil {
		t.Fatalf("list submissions by user: %v", err)
	}
	if total != 1 || len(submissions) != 1 || submissions[0].ID != created.ID {
		t.Fatalf("unexpected submission list: total=%d submissions=%#v", total, submissions)
	}
}

func TestCreateSubmissionWithAnswersRollsBackTogether(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Submission{}, &Answer{}); err != nil {
		t.Fatalf("migrate submission models: %v", err)
	}
	user := User{Name: "Learner", Email: "learner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Exec(`CREATE TRIGGER reject_answer BEFORE INSERT ON answer BEGIN SELECT RAISE(FAIL, 'answer rejected'); END`).Error; err != nil {
		t.Fatalf("create answer trigger: %v", err)
	}

	submission := &Submission{UserID: user.ID, QuestionSetID: 8, AnswersCount: 1}
	_, err = NewRepository(db).CreateSubmissionWithAnswers(submission, []answer.Answer{{
		QuestionID: 12, QuestionOptionID: 34, IsCorrect: true,
	}})
	if err == nil {
		t.Fatal("expected answer persistence to fail")
	}
	var count int64
	if err := db.Model(&Submission{}).Count(&count).Error; err != nil {
		t.Fatalf("count submissions: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected submission to roll back with failed answer, got %d rows", count)
	}
}
