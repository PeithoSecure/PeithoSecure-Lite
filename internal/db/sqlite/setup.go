package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the SQLite connection using env or fallback path
func InitDB() {
	dbPath := os.Getenv("PEITHO_SQLITE_PATH")
	if dbPath == "" {
		dbPath = "peitho_secure.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}

	// Enforce WAL mode for better concurrency
	if _, err := DB.Exec(`PRAGMA journal_mode = WAL;`); err != nil {
		log.Fatalf("❌ Failed to set WAL mode: %v", err)
	}

	createTables()
	verifySchema()
}

// createTables defines all required tables for PeithoSecure Lite
func createTables() {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			role TEXT NOT NULL,
			email_verified INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS roast_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			event TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS email_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS audit_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			event_type TEXT NOT NULL,
			ip_address TEXT,
			user_agent TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, stmt := range stmts {
		if _, err := DB.Exec(stmt); err != nil {
			log.Fatalf("❌ Failed to execute table init: %v\nSQL: %s", err, stmt)
		}
	}

	log.Println("✅ Database tables initialized successfully.")
}

// verifySchema performs a runtime sanity check for table presence
func verifySchema() {
	requiredTables := []string{
		"users",
		"roast_logs",
		"email_tokens",
		"audit_events",
	}

	for _, table := range requiredTables {
		query := fmt.Sprintf(`SELECT name FROM sqlite_master WHERE type='table' AND name='%s';`, table)
		var name string
		if err := DB.QueryRow(query).Scan(&name); err != nil {
			log.Fatalf("🚨 Missing required table: %s. Please reset your database.", table)
		}
	}
	log.Println("🧪 Schema verification passed: All required tables exist.")
}

// GetDB exposes the active SQLite DB instance
func GetDB() *sql.DB {
	return DB
}
