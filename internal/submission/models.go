package submission

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint `gorm:"not null"`
	QuestionSetID uint `gorm:"not null"`
	CorrectCount  int  `gorm:"not null"`
	AnswersCount  int  `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	User          *User          `gorm:"foreignKey:UserID"`
	Answers       []Answer       `gorm:"foreignKey:SubmissionID"`
}

func (Submission) TableName() string {
	return "submission"
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

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
