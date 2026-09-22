package answer

import (
	"time"

	"gorm.io/gorm"
)

// GORM persistence models owned by the answer domain. None of them carry json
// tags: serialization belongs to the DTOs in internal/answer/dto.

type Answer struct {
	ID               uint `gorm:"primaryKey"`
	SubmissionID     uint `gorm:"not null"`
	QuestionID       uint `gorm:"not null"`
	QuestionOptionID uint `gorm:"not null"`
	IsCorrect        bool `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (Answer) TableName() string {
	return "answer"
}

// Subject is the answer view of the subject table, limited to what the
// performance projection embeds.
type Subject struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func (Subject) TableName() string {
	return "subject"
}
