package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/direwen/go-server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func ConnectDB() {
	var dsn string

	// Build dsn
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

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	sqlDB.SetMaxOpenConns(25)                 // max connections open at once
	sqlDB.SetMaxIdleConns(10)                 // connections kept idle in pool
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // how long a connection can be reused

	log.Println("Connected to Database")

	// Migrate Database
	err = DB.AutoMigrate(
		&model.ContextTemplate{},
		&model.Scenario{},
		&model.Session{},
		&model.Response{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
	log.Println("Database migrated")
}
