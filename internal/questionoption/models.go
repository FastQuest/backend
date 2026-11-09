package questionoption

type QuestionOption struct {
	ID         uint   `gorm:"primaryKey"`
	Text       string `gorm:"not null"`
	IsCorrect  bool   `gorm:"not null"`
	QuestionID uint   `gorm:"column:id_question;not null"`
}

func (QuestionOption) TableName() string {
	return "question_option"
}
