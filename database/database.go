package database

import (
	internaldb "flashquest/internal/platform/database"

	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	return internaldb.InitDB()
}

func GetDB() *gorm.DB {
	return internaldb.GetDB()
}
