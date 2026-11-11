package questionset

import "time"

type QuestionSet struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	Type        string `gorm:"not null"`
	UserID      int
	User        *User `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	IsPrivate   bool `gorm:"not null"`
}

func (QuestionSet) TableName() string {
	return "question_set"
}

type QuestionSetQuestion struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	QuestionSetID int
	QuestionID    int
	Position      int
}

func (QuestionSetQuestion) TableName() string {
	return "question_set_question"
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
