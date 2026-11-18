package classroom

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testQuestionSet struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	IsPrivate bool
}

func (testQuestionSet) TableName() string { return "question_set" }

type testQuestionSetQuestion struct {
	ID            uint `gorm:"primaryKey"`
	QuestionSetID uint
	QuestionID    uint
	Position      int
}

func (testQuestionSetQuestion) TableName() string { return "question_set_question" }

type testQuestion struct {
	ID        uint `gorm:"primaryKey"`
	Statement string
}

func (testQuestion) TableName() string { return "question" }

type testQuestionOption struct {
	ID         uint `gorm:"primaryKey"`
	QuestionID uint `gorm:"column:id_question"`
	Text       string
	IsCorrect  bool
}

func (testQuestionOption) TableName() string { return "question_option" }

type testUser struct {
	ID uint `gorm:"primaryKey"`
}

func (testUser) TableName() string { return "users" }

func openClassroomTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(
		&testUser{}, &Classroom{}, &Participant{}, &QuestionSetVersion{},
		&VersionQuestion{}, &VersionOption{}, &Activity{}, &testQuestionSet{},
		&testQuestionSetQuestion{}, &testQuestion{}, &testQuestionOption{},
	); err != nil {
		t.Fatalf("migrate test models: %v", err)
	}
	return db
}

func TestCreateActivityFreezesQuestionAndOptionContent(t *testing.T) {
	db := openClassroomTestDB(t)
	user := testUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	classroom := Classroom{Name: "Biology", CreatorID: user.ID, Status: StatusActive}
	if err := db.Create(&classroom).Error; err != nil {
		t.Fatalf("create classroom: %v", err)
	}
	repo := NewRepository(db)
	if _, err := repo.AddParticipant(classroom.ID, user.ID); err != nil {
		t.Fatalf("add participant: %v", err)
	}
	questionSet := testQuestionSet{UserID: user.ID}
	if err := db.Create(&questionSet).Error; err != nil {
		t.Fatalf("create question set: %v", err)
	}
	question := testQuestion{Statement: "Original question"}
	if err := db.Create(&question).Error; err != nil {
		t.Fatalf("create question: %v", err)
	}
	option := testQuestionOption{QuestionID: question.ID, Text: "Original option", IsCorrect: true}
	if err := db.Create(&option).Error; err != nil {
		t.Fatalf("create question option: %v", err)
	}
	if err := db.Create(&testQuestionSetQuestion{QuestionSetID: questionSet.ID, QuestionID: question.ID, Position: 1}).Error; err != nil {
		t.Fatalf("add question to set: %v", err)
	}

	activity, err := repo.CreateActivity(classroom.ID, questionSet.ID, user.ID, "Quiz", false)
	if err != nil {
		t.Fatalf("create activity: %v", err)
	}
	if err := db.Model(&testQuestion{}).Where("id = ?", question.ID).Update("statement", "Changed question").Error; err != nil {
		t.Fatalf("change source question: %v", err)
	}
	if err := db.Model(&testQuestionOption{}).Where("id = ?", option.ID).Updates(map[string]any{"text": "Changed option", "is_correct": false}).Error; err != nil {
		t.Fatalf("change source option: %v", err)
	}

	context, err := repo.GetSubmissionContext(activity.ID, user.ID)
	if err != nil {
		t.Fatalf("get activity submission context: %v", err)
	}
	if context.QuestionSetID != questionSet.ID || len(context.Questions) != 1 || context.Questions[0].ID != question.ID {
		t.Fatalf("unexpected snapshot: %+v", context)
	}
	correct, ok := context.Questions[0].Options[option.ID]
	if !ok || !correct {
		t.Fatalf("expected original option correctness to remain true, got options %#v", context.Questions[0].Options)
	}
	var frozenQuestion VersionQuestion
	if err := db.First(&frozenQuestion, "version_id = ?", activity.VersionID).Error; err != nil {
		t.Fatalf("load frozen question: %v", err)
	}
	if frozenQuestion.Statement != "Original question" {
		t.Fatalf("expected original question text, got %q", frozenQuestion.Statement)
	}
	var frozenOption VersionOption
	if err := db.First(&frozenOption, "version_question_id = ?", frozenQuestion.ID).Error; err != nil {
		t.Fatalf("load frozen option: %v", err)
	}
	if frozenOption.Text != "Original option" {
		t.Fatalf("expected original option text, got %q", frozenOption.Text)
	}
	activities, err := NewService(repo, db).ListActivities(user.ID, "Aluno", classroom.ID)
	if err != nil {
		t.Fatalf("list activities for participant: %v", err)
	}
	if len(activities) != 1 || len(activities[0].Questions) != 1 || activities[0].Questions[0].Statement != "Original question" {
		t.Fatalf("expected listed activity to use the frozen question, got %#v", activities)
	}
	payload, err := json.Marshal(activities)
	if err != nil {
		t.Fatalf("marshal activity response: %v", err)
	}
	if strings.Contains(string(payload), "is_correct") {
		t.Fatalf("activity response must not expose the answer key: %s", payload)
	}
	if _, err := repo.GetSubmissionContext(activity.ID, user.ID+100); !errors.Is(err, ErrNotActiveParticipant) {
		t.Fatalf("expected nonparticipant to be rejected, got %v", err)
	}
}

func TestCreateActivityRollsBackWhenListCannotBeSnapshotted(t *testing.T) {
	db := openClassroomTestDB(t)
	user := testUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	classroom := Classroom{Name: "Empty", CreatorID: user.ID, Status: StatusActive}
	if err := db.Create(&classroom).Error; err != nil {
		t.Fatalf("create classroom: %v", err)
	}
	questionSet := testQuestionSet{UserID: user.ID}
	if err := db.Create(&questionSet).Error; err != nil {
		t.Fatalf("create question set: %v", err)
	}

	_, err := NewRepository(db).CreateActivity(classroom.ID, questionSet.ID, user.ID, "Empty activity", false)
	if !errors.Is(err, ErrQuestionSetEmpty) {
		t.Fatalf("expected empty question set error, got %v", err)
	}
	var versions, activities int64
	if err := db.Model(&QuestionSetVersion{}).Count(&versions).Error; err != nil {
		t.Fatalf("count versions: %v", err)
	}
	if err := db.Model(&Activity{}).Count(&activities).Error; err != nil {
		t.Fatalf("count activities: %v", err)
	}
	if versions != 0 || activities != 0 {
		t.Fatalf("expected failed activity creation to roll back, got versions=%d activities=%d", versions, activities)
	}
}
