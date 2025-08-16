package database

import (
	"donkey/config"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

// DB instance
var (
	DB  *gorm.DB
	err error
)

// Connect to the database
func Connect() {

	dsn := ""

	switch os.Getenv("DB_KIND") {
	case "postgreslq", "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		// github.com/mattn/go-sqlite3
		DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", os.Getenv("DB_NAME"))), &gorm.Config{})
	default:
		slog.Error("Fatal: the database used isn't supported. Must be one of `sqlite`, `postgres`", slog.Any("err", err))
		panic("Incorrect database type used")
	}

	if err != nil {
		slog.Error("Fatal: faile to connect to the database!", slog.Any("err", err))
		panic("Failed to connect database")
	}
	slog.Info("Connection Opened to Database")

	// Migrate the schemas
	DB.AutoMigrate(&config.Book{}, &config.User{})
	slog.Info("Database Migrated")
}
