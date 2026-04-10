package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	env := os.Getenv("RAILWAY_ENVIRONMENT")
	if env == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("Aviso: .env não carregado")
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	log.Println(dsn)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if rawDB, err := db.DB(); err == nil {
		rawDB.Exec("DEALLOCATE ALL")
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get generic DB interface:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	DB = db
	log.Println("Database connection established successfully")
	return db
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection not initialized. Call InitDB() first.")
	}
	return DB
}
