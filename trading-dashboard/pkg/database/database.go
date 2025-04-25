package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var DB *sql.DB

// Initialize sets up the database connection
func Initialize() error {
	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "trading.db")
	log.Printf("Initializing database at: %s", dbPath)

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize schema
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// createTables creates the database tables if they don't exist
func createTables() error {
	// Risk Assessment table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS risk_assessments (
		id INTEGER PRIMARY KEY,
		date TEXT,
		emotional INTEGER,
		fomo INTEGER,
		bias INTEGER,
		physical INTEGER,
		pnl INTEGER,
		overall_score INTEGER
	)`)
	if err != nil {
		return fmt.Errorf("failed to create risk_assessments table: %w", err)
	}

	// Stock Ratings table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS stock_ratings (
		id INTEGER PRIMARY KEY,
		date TEXT,
		ticker TEXT,
		market_sentiment INTEGER,
		basic_materials INTEGER,
		communication_services INTEGER,
		consumer_cyclical INTEGER,
		consumer_defensive INTEGER,
		energy INTEGER,
		financial INTEGER,
		healthcare INTEGER,
		industrials INTEGER,
		real_estate INTEGER,
		technology INTEGER,
		utilities INTEGER,
		stock_sentiment INTEGER,
		pattern TEXT,
		enthusiasm_rating INTEGER
	)`)
	if err != nil {
		return fmt.Errorf("failed to create stock_ratings table: %w", err)
	}

	// Trades table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY,
		entry_date TEXT,
		ticker TEXT,
		sector TEXT,
		entry_price REAL,
		notes TEXT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create trades table: %w", err)
	}

	return nil
}
