package models

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	SubmissionID     uint           `gorm:"not null" json:"submission_id"`
	QuestionID       uint           `gorm:"not null" json:"question_id"`
	QuestionOptionID uint           `gorm:"not null" json:"option_id"`
	IsCorrect        bool           `gorm:"not null" json:"is_correct"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Answer) TableName() string {
	return "answer"
}
