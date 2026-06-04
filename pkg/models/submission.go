package models

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	QuestionSetID uint           `gorm:"not null" json:"question_set_id"`
	CorrectCount  int            `gorm:"not null" json:"correct_count"`
	AnswersCount  int            `gorm:"not null" json:"answers_count"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Submission) TableName() string {
	return "submission"
}
