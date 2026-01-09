package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/direwen/go-server/internal/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func ConnectDB() {
	var dsn string

	if os.Getenv("DB_URL") != "" {
		dsn = os.Getenv("DB_URL")
	} else {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to Database")

	err = DB.AutoMigrate(
		&models.ContextTemplate{},
		&models.Session{},
		&models.Scenario{},
		&models.Response{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
	log.Println("Database migrated")
}
