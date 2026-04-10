package source

import (
	"flashquest/database"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}
