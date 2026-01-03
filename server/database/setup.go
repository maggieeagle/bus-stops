package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func SetDB(pgdb *sql.DB) *sql.DB {
	database_url := os.Getenv("DATABASE_URL")

	// PostgreSQL
	pgDSN := fmt.Sprintf("%s", database_url)
	var err error
	pgdb, err = sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	pgdb.SetConnMaxLifetime(time.Minute * 3)
	pgdb.SetMaxOpenConns(10)
	pgdb.SetMaxIdleConns(10)
	if err = pgdb.Ping(); err != nil {
		log.Fatalf("PostgreSQL ping failed: %v", err)
	}

	return pgdb
}