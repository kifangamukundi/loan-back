package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB establishes a connection to the database and sets the global DB variable.
func ConnectDB() error {
	dbType := os.Getenv("DB_TYPE") // Check which database to use

	var err error
	switch dbType {
	case "postgres":
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			return fmt.Errorf("DATABASE_URL not set in environment")
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	case "sqlite":
		dbFile := os.Getenv("SQLITE_FILE")
		if dbFile == "" {
			dbFile = "database.db" // Default to `database.db` if no file is provided
		}
		DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

	default:
		return fmt.Errorf("Invalid DB_TYPE. Use 'postgres' or 'sqlite'")
	}

	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB returns the initialized database connection instance.
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection is not initialized. Call ConnectDB first.")
	}
	return DB
}
