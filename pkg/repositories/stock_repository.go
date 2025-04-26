package repositories

import (
	"fmt"
	"strings"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

const STOCK_PREFIX = "stock_"

// SaveStockRating saves a stock rating to the database
func SaveStockRating(rating *models.StockRating) error {
	// Calculate enthusiasm rating
	rating.CalculateEnthusiasm()

	// If no ID is set, generate one
	if rating.ID == "" {
		// Generate key with ticker included for easier filtering
		rating.ID = fmt.Sprintf("%s%s_%s", STOCK_PREFIX, rating.Ticker, rating.Date.Format("20060102"))
	}

	// Save the rating to BadgerDB
	return database.Set(rating.ID, rating)
}

// GetStockRating retrieves a stock rating by ID
func GetStockRating(id string) (*models.StockRating, error) {
	rating := &models.StockRating{}
	err := database.Get(id, rating)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock rating: %w", err)
	}
	return rating, nil
}

// GetStockRatingsByTicker retrieves all stock ratings for a specific ticker
func GetStockRatingsByTicker(ticker string) ([]*models.StockRating, error) {
	// Get all stock ratings
	allRatings, err := GetAllStockRatings()
	if err != nil {
		return nil, err
	}

	// Filter ratings by ticker
	var filteredRatings []*models.StockRating
	for _, rating := range allRatings {
		if strings.EqualFold(rating.Ticker, ticker) {
			filteredRatings = append(filteredRatings, rating)
		}
	}

	return filteredRatings, nil
}

// GetAllStockRatings retrieves all stock ratings
func GetAllStockRatings() ([]*models.StockRating, error) {
	var ratings []*models.StockRating
	err := database.GetByPrefix(STOCK_PREFIX, &ratings)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock ratings: %w", err)
	}
	return ratings, nil
}

// DeleteStockRating deletes a stock rating by ID
func DeleteStockRating(id string) error {
	return database.Delete(id)
}
