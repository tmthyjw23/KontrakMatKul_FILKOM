package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env") // Load from backend root

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	query := `DELETE FROM registrations WHERE student_nim = '22010001'`
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalf("exec: %v", err)
	}

	rows, _ := res.RowsAffected()
	log.Printf("Successfully deleted %d registrations for Budi Santoso (22010001)\n", rows)
}
