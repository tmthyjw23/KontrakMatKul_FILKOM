// cmd/seeder/main.go
//
// Run:  go run ./cmd/seeder/main.go
//
// This program reads docs/migrations/003_seed_data.sql and executes it
// against the DATABASE_URL from .env (or the environment).
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// ── 1. Load .env ─────────────────────────────────────────
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌  DATABASE_URL is not set")
	}

	// ── 2. Connect ───────────────────────────────────────────
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌  open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌  ping: %v", err)
	}
	log.Println("✅  Connected to database")

	// ── 3. Run migration 002 (schema) ────────────────────────
	runFile(db, "docs/migrations/002_add_lecturers_prerequisites_curriculums.sql")

	// ── 4. Run seed 003 (data) ───────────────────────────────
	runFile(db, "docs/migrations/003_seed_data.sql")

	log.Println("🎉  Seeding complete!")
}

// runFile executes every SQL statement in a file.
// It splits on ';' and skips blank / comment-only statements.
func runFile(db *sql.DB, path string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("❌  read %s: %v", path, err)
	}

	// Remove comment lines before splitting by semicolon
	var cleanLines []string
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue // skip comment line
		}
		cleanLines = append(cleanLines, line)
	}

	cleanSQL := strings.Join(cleanLines, "\n")
	stmts := strings.Split(cleanSQL, ";")

	count := 0
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			// Print the failing statement (first 120 chars) for diagnostics
			preview := stmt
			if len(preview) > 120 {
				preview = preview[:120] + "…"
			}
			log.Printf("⚠️   skipped (%.120s…): %v", preview, err)
			continue
		}
		count++
	}
	log.Printf("✅  %s — executed %d statements", path, count)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func mustExec(db *sql.DB, query string) sql.Result {
	res, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌  exec: %v\nquery: %s", err, query)
	}
	return res
}

func check(err error) {
	if err != nil {
		log.Fatal("❌ ", err)
	}
}

func init() {
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("  KontrakMatKul FILKOM — Database Seeder  ")
	fmt.Println("──────────────────────────────────────────")
}
