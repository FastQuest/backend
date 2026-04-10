package exam

import (
	"errors"
	"flashquest/database"
	"flashquest/pkg/models"

	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}

func GetInstanceWithSource(eiID int, ei *models.ExamInstance) error {
	db := getDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	result := db.Preload("SourceExamInstance.Source").Where("id = ?", eiID).Find(&ei)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Source Exam Instance not found")
		}
		return errors.New("Error fetching Source Exam Instance")
	}

	return nil
}
