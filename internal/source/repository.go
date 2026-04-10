package source

import (
	database "flashquest/internal/platform/database"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return database.GetDB()
}
