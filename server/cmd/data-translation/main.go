package main

import (
	"encoding/csv"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func copyTable(db *sql.DB, table string, file string) {
	f, err := os.Open(filepath.Join("data", file))
	if err != nil {
		log.Fatalf("open %s: %v", file, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// read column headers
	columns, err := reader.Read()
	if err != nil {
		log.Fatalf("read header %s: %v", file, err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(pq.CopyIn(table, columns...))
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		values := make([]any, len(record))
		for i, v := range record {
			if v == "" {
				values[i] = nil
			} else {
				values[i] = v
			}
		}

		if _, err := stmt.Exec(values...); err != nil {
			log.Fatalf("copy %s: %v", table, err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}

	if err := stmt.Close(); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Imported %s\n", table)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	copyTable(db, "routes", "routes.txt")
	copyTable(db, "stops", "stops.txt")
	copyTable(db, "trips", "trips.txt")
	copyTable(db, "stop_times", "stop_times.txt")

	log.Println("All data imported successfully")
}
