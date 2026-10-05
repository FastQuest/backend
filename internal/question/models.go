package question

import "time"

// GORM persistence models the question entity preloads. They are the question
// domain view of tables owned by other contexts, limited to what the question
// responses expose, and none of them carry json tags.

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func (Subject) TableName() string {
	return "subject"
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

type QuestionOption struct {
	ID         uint   `gorm:"primaryKey"`
	Text       string `gorm:"not null"`
	IsCorrect  bool   `gorm:"not null"`
	QuestionID uint   `gorm:"column:id_question; not null"`
}

func (QuestionOption) TableName() string {
	return "question_option"
}

type ExamInstance struct {
	ID        uint   `gorm:"primaryKey"`
	SourceId  uint   `gorm:"not null"`
	Source    Source `gorm:"foreignKey:SourceId"`
	Edition   uint   `gorm:"not null"`
	Phase     uint   `gorm:"not null"`
	Year      uint   `gorm:"not null"`
	CreatedAt time.Time
}

func (ExamInstance) TableName() string {
	return "source_exam_instance"
}

type Source struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	CreatedAt time.Time
}

func (Source) TableName() string {
	return "source"
}
