package models

type QuestionOption struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Text       string `gorm:"not null" json:"text"`
	Is_correct bool   `gorm:"not null" json:"is_correct"`
	QuestionID uint   `gorm:"column:id_question; not null" json:"question_id"`
}

func (QuestionOption) TableName() string {
	return "question_option"
}
