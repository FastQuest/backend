package classroom

import "time"

const (
	StatusActive   = "active"
	StatusArchived = "archived"
)

type Classroom struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	CreatorID uint   `gorm:"not null"`
	Status    string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Classroom) TableName() string { return "classroom" }

type Participant struct {
	ID          uint   `gorm:"primaryKey"`
	ClassroomID uint   `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	Status      string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Participant) TableName() string { return "classroom_participant" }

type QuestionSetVersion struct {
	ID            uint `gorm:"primaryKey"`
	QuestionSetID uint `gorm:"not null"`
	VersionNumber int  `gorm:"not null"`
	CreatedBy     uint `gorm:"not null"`
	CreatedAt     time.Time
	Questions     []VersionQuestion `gorm:"foreignKey:VersionID"`
}

func (QuestionSetVersion) TableName() string { return "question_set_version" }

type VersionQuestion struct {
	ID               uint            `gorm:"primaryKey"`
	VersionID        uint            `gorm:"not null"`
	SourceQuestionID uint            `gorm:"not null"`
	Position         int             `gorm:"not null"`
	Statement        string          `gorm:"not null"`
	Options          []VersionOption `gorm:"foreignKey:VersionQuestionID"`
}

func (VersionQuestion) TableName() string { return "question_set_version_question" }

type VersionOption struct {
	ID                uint   `gorm:"primaryKey"`
	VersionQuestionID uint   `gorm:"not null"`
	SourceOptionID    uint   `gorm:"not null"`
	Position          int    `gorm:"not null"`
	Text              string `gorm:"not null"`
	IsCorrect         bool   `gorm:"not null"`
}

func (VersionOption) TableName() string { return "question_set_version_option" }

type Activity struct {
	ID          uint   `gorm:"primaryKey"`
	ClassroomID uint   `gorm:"not null"`
	VersionID   uint   `gorm:"not null"`
	Title       string `gorm:"not null"`
	Status      string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     QuestionSetVersion `gorm:"foreignKey:VersionID"`
}

func (Activity) TableName() string { return "classroom_activity" }
