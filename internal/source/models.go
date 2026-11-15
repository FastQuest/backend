package source

import "time"

type Source struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	CreatedAt time.Time
}

func (Source) TableName() string {
	return "source"
}

type ExamInstance struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SourceID  uint   `gorm:"not null" json:"source_id"`
	Source    Source `gorm:"foreignKey:SourceID" json:"source"`
	Edition   uint   `gorm:"not null" json:"edition"`
	Phase     uint   `gorm:"not null" json:"phase"`
	Year      uint   `gorm:"not null" json:"year"`
	CreatedAt time.Time
}

func (ExamInstance) TableName() string {
	return "source_exam_instance"
}
