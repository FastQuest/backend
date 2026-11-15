package source

import (
	"encoding/json"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestExamInstancePreservesJSONContract(t *testing.T) {
	body, err := json.Marshal(ExamInstance{
		ID: 1,
		Source: Source{
			ID:   2,
			Name: "Example",
			Type: "exam",
		},
		SourceID: 2,
		Edition:  3,
		Phase:    4,
		Year:     2024,
	})
	if err != nil {
		t.Fatalf("marshal exam instance: %v", err)
	}

	want := `{"id":1,"source_id":2,"source":{"ID":2,"Name":"Example","Type":"exam","CreatedAt":"0001-01-01T00:00:00Z"},"edition":3,"phase":4,"year":2024,"CreatedAt":"0001-01-01T00:00:00Z"}`
	if string(body) != want {
		t.Fatalf("expected JSON %s, got %s", want, body)
	}
}

func TestCreateSourcePersistsExamInstanceWithSourceIDColumn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(&Source{}, &ExamInstance{}); err != nil {
		t.Fatalf("migrate source models: %v", err)
	}

	repository := NewRepository(db)
	sourceID, err := repository.CreateSource(SourceExamBody{
		Name:    "Example",
		Type:    "exam",
		Edition: 2,
		Phase:   1,
		Year:    2024,
	})
	if err != nil {
		t.Fatalf("create source: %v", err)
	}

	var instance ExamInstance
	if err := db.Where("source_id = ?", sourceID).First(&instance).Error; err != nil {
		t.Fatalf("find persisted exam instance: %v", err)
	}
	if instance.SourceID != sourceID || instance.Edition != 2 || instance.Phase != 1 || instance.Year != 2024 {
		t.Fatalf("unexpected exam instance: %+v", instance)
	}
}
