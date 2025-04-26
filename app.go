package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
	"trading-dashboard/pkg/repositories"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Set up logging to a file
	setupLogging()

	log.Println("===== Application Starting =====")
	return &App{}
}

// setupLogging configures logging to write to both console and file
func setupLogging() {
	// Create logs directory in user config dir
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Warning: Could not get user config directory: %v", err)
		return
	}

	logDir := filepath.Join(appDataDir, "TradingDashboard", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Warning: Could not create log directory: %v", err)
		return
	}

	// Create log file with timestamp
	logFile := filepath.Join(logDir, "app-"+time.Now().Format("20060102-150405")+".log")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Warning: Could not create log file: %v", err)
		return
	}

	// Set up multi-writer to log to both console and file
	log.SetOutput(f)
	log.Printf("Log file created at: %s", logFile)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("Starting application in context:", ctx)

	// Initialize database
	log.Println("Initializing database...")
	if err := database.Initialize(); err != nil {
		log.Fatalf("FATAL: Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Log app data and database path
	appDataDir, err := os.UserConfigDir()
	if err == nil {
		dataDir := filepath.Join(appDataDir, "TradingDashboard", "data")
		log.Printf("Using data directory: %s", dataDir)

		// Check if directory exists and log content
		files, err := os.ReadDir(dataDir)
		if err != nil {
			log.Printf("Error reading data directory: %v", err)
		} else {
			log.Printf("Found %d items in data directory", len(files))
			for _, file := range files {
				log.Printf("  - %s (%d bytes)", file.Name(), getFileSize(filepath.Join(dataDir, file.Name())))
			}
		}
	}
}

// Helper to get file size
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	log.Println("Application shutting down...")
	// Close database connection
	database.Close()
	log.Println("Database connection closed")
	log.Println("===== Application Terminated =====")
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return "0.1.0"
}

// Risk Management API Methods

// SaveRiskAssessment saves a risk assessment
func (a *App) SaveRiskAssessment(assessment models.RiskAssessment) (*models.RiskAssessment, error) {
	log.Printf("API: SaveRiskAssessment called with ID=%s", assessment.ID)
	err := repositories.SaveRiskAssessment(&assessment)
	if err != nil {
		log.Printf("ERROR: SaveRiskAssessment failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: SaveRiskAssessment saved with ID=%s", assessment.ID)
	return &assessment, nil
}

// GetLatestRiskAssessment gets the latest risk assessment
func (a *App) GetLatestRiskAssessment() (*models.RiskAssessment, error) {
	log.Println("API: GetLatestRiskAssessment called")
	result, err := repositories.GetLatestRiskAssessment()
	if err != nil {
		log.Printf("ERROR: GetLatestRiskAssessment failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: GetLatestRiskAssessment returned ID=%s", result.ID)
	return result, nil
}

// GetAllRiskAssessments gets all risk assessments
func (a *App) GetAllRiskAssessments() ([]*models.RiskAssessment, error) {
	log.Println("API: GetAllRiskAssessments called")
	result, err := repositories.GetAllRiskAssessments()
	if err != nil {
		log.Printf("ERROR: GetAllRiskAssessments failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: GetAllRiskAssessments returned %d records", len(result))
	return result, nil
}

// Stock Rating API Methods

// SaveStockRating saves a stock rating
func (a *App) SaveStockRating(rating models.StockRating) (*models.StockRating, error) {
	log.Printf("API: SaveStockRating called with ticker=%s", rating.Ticker)
	err := repositories.SaveStockRating(&rating)
	if err != nil {
		log.Printf("ERROR: SaveStockRating failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: SaveStockRating saved with ID=%s", rating.ID)
	return &rating, nil
}

// GetStockRating gets a stock rating by ID
func (a *App) GetStockRating(id string) (*models.StockRating, error) {
	log.Printf("API: GetStockRating called with ID=%s", id)
	return repositories.GetStockRating(id)
}

// GetStockRatingsByTicker gets all stock ratings for a specific ticker
func (a *App) GetStockRatingsByTicker(ticker string) ([]*models.StockRating, error) {
	log.Printf("API: GetStockRatingsByTicker called with ticker=%s", ticker)
	return repositories.GetStockRatingsByTicker(ticker)
}

// GetAllStockRatings gets all stock ratings
func (a *App) GetAllStockRatings() ([]*models.StockRating, error) {
	log.Println("API: GetAllStockRatings called")
	result, err := repositories.GetAllStockRatings()
	if err != nil {
		log.Printf("ERROR: GetAllStockRatings failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: GetAllStockRatings returned %d records", len(result))
	return result, nil
}

// Trade Calendar API Methods

// SaveTrade saves a trade
func (a *App) SaveTrade(trade models.Trade) (*models.Trade, error) {
	log.Printf("API: SaveTrade called with ticker=%s", trade.Ticker)
	err := repositories.SaveTrade(&trade)
	if err != nil {
		log.Printf("ERROR: SaveTrade failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: SaveTrade saved with ID=%s", trade.ID)
	return &trade, nil
}

// GetTrade gets a trade by ID
func (a *App) GetTrade(id string) (*models.Trade, error) {
	log.Printf("API: GetTrade called with ID=%s", id)
	return repositories.GetTrade(id)
}

// GetTradesByDateRange gets trades within a date range
func (a *App) GetTradesByDateRange(startDateStr, endDateStr string) ([]*models.Trade, error) {
	log.Printf("API: GetTradesByDateRange called with range=%s to %s", startDateStr, endDateStr)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Printf("ERROR: Failed to parse start date: %v", err)
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Printf("ERROR: Failed to parse end date: %v", err)
		return nil, err
	}

	result, err := repositories.GetTradesByDateRange(startDate, endDate)
	if err != nil {
		log.Printf("ERROR: GetTradesByDateRange failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: GetTradesByDateRange returned %d records", len(result))
	return result, nil
}

// GetTradesByTicker gets trades for a specific ticker
func (a *App) GetTradesByTicker(ticker string) ([]*models.Trade, error) {
	log.Printf("API: GetTradesByTicker called with ticker=%s", ticker)
	return repositories.GetTradesByTicker(ticker)
}

// GetAllTrades gets all trades
func (a *App) GetAllTrades() ([]*models.Trade, error) {
	log.Println("API: GetAllTrades called")
	result, err := repositories.GetAllTrades()
	if err != nil {
		log.Printf("ERROR: GetAllTrades failed: %v", err)
		return nil, err
	}
	log.Printf("SUCCESS: GetAllTrades returned %d records", len(result))
	return result, nil
}

// DeleteTrade deletes a trade
func (a *App) DeleteTrade(id string) error {
	log.Printf("API: DeleteTrade called with ID=%s", id)
	err := repositories.DeleteTrade(id)
	if err != nil {
		log.Printf("ERROR: DeleteTrade failed: %v", err)
		return err
	}
	log.Printf("SUCCESS: DeleteTrade completed for ID=%s", id)
	return nil
}
