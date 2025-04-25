package main

import (
	"context"
	"log"
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
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Close database connection
	database.Close()
	log.Println("Database connection closed")
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return "0.1.0"
}

// Risk Management API Methods

// SaveRiskAssessment saves a risk assessment
func (a *App) SaveRiskAssessment(assessment models.RiskAssessment) (*models.RiskAssessment, error) {
	err := repositories.SaveRiskAssessment(&assessment)
	if err != nil {
		return nil, err
	}
	return &assessment, nil
}

// GetLatestRiskAssessment gets the latest risk assessment
func (a *App) GetLatestRiskAssessment() (*models.RiskAssessment, error) {
	return repositories.GetLatestRiskAssessment()
}

// GetAllRiskAssessments gets all risk assessments
func (a *App) GetAllRiskAssessments() ([]*models.RiskAssessment, error) {
	return repositories.GetAllRiskAssessments()
}

// Stock Rating API Methods

// SaveStockRating saves a stock rating
func (a *App) SaveStockRating(rating models.StockRating) (*models.StockRating, error) {
	err := repositories.SaveStockRating(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetStockRating gets a stock rating by ID
func (a *App) GetStockRating(id int) (*models.StockRating, error) {
	return repositories.GetStockRating(id)
}

// GetStockRatingsByTicker gets all stock ratings for a specific ticker
func (a *App) GetStockRatingsByTicker(ticker string) ([]*models.StockRating, error) {
	return repositories.GetStockRatingsByTicker(ticker)
}

// GetAllStockRatings gets all stock ratings
func (a *App) GetAllStockRatings() ([]*models.StockRating, error) {
	return repositories.GetAllStockRatings()
}

// Trade Calendar API Methods

// SaveTrade saves a trade
func (a *App) SaveTrade(trade models.Trade) (*models.Trade, error) {
	err := repositories.SaveTrade(&trade)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

// GetTrade gets a trade by ID
func (a *App) GetTrade(id int) (*models.Trade, error) {
	return repositories.GetTrade(id)
}

// GetTradesByDateRange gets trades within a date range
func (a *App) GetTradesByDateRange(startDateStr, endDateStr string) ([]*models.Trade, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, err
	}

	return repositories.GetTradesByDateRange(startDate, endDate)
}

// GetTradesByTicker gets trades for a specific ticker
func (a *App) GetTradesByTicker(ticker string) ([]*models.Trade, error) {
	return repositories.GetTradesByTicker(ticker)
}

// GetAllTrades gets all trades
func (a *App) GetAllTrades() ([]*models.Trade, error) {
	return repositories.GetAllTrades()
}
