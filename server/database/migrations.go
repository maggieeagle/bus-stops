package database

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations() {
	migrationsPath := os.Getenv("MIGRATIONS_FOLDER_PATH")
	dbURL := os.Getenv("DATABASE_URL")

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	err = m.Up()
	switch err {
	case nil:
		log.Println("Migrations applied successfully.")
	case migrate.ErrNoChange:
		log.Println("No new migrations to apply.")
	default:
		log.Fatalf("Migration failed: %v", err)
	}

	// Close the migrate tools DB connections
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		log.Printf("Warning: source close error: %v", sourceErr)
	}
	if dbErr != nil {
		log.Printf("Warning: database close error: %v", dbErr)
	}
}