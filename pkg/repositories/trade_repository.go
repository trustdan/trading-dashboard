package repositories

import (
	"fmt"
	"strings"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

const TRADE_PREFIX = "trade_"

// SaveTrade saves a trade to the database
func SaveTrade(trade *models.Trade) error {
	// If no ID is set, generate one
	if trade.ID == "" {
		trade.ID = database.GenerateKey(TRADE_PREFIX)
	}

	// Save the trade to BadgerDB
	return database.Set(trade.ID, trade)
}

// GetTrade retrieves a trade by ID
func GetTrade(id string) (*models.Trade, error) {
	trade := &models.Trade{}
	err := database.Get(id, trade)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}
	return trade, nil
}

// GetTradesByDateRange retrieves trades within a date range
func GetTradesByDateRange(startDate, endDate time.Time) ([]*models.Trade, error) {
	// Get all trades
	allTrades, err := GetAllTrades()
	if err != nil {
		return nil, err
	}

	// Filter trades by date range
	var filteredTrades []*models.Trade
	for _, trade := range allTrades {
		if (trade.EntryDate.Equal(startDate) || trade.EntryDate.After(startDate)) &&
			(trade.EntryDate.Equal(endDate) || trade.EntryDate.Before(endDate)) {
			filteredTrades = append(filteredTrades, trade)
		}
	}

	return filteredTrades, nil
}

// GetTradesByTicker retrieves trades for a specific ticker
func GetTradesByTicker(ticker string) ([]*models.Trade, error) {
	// Get all trades
	allTrades, err := GetAllTrades()
	if err != nil {
		return nil, err
	}

	// Filter trades by ticker
	var filteredTrades []*models.Trade
	for _, trade := range allTrades {
		if strings.EqualFold(trade.Ticker, ticker) {
			filteredTrades = append(filteredTrades, trade)
		}
	}

	return filteredTrades, nil
}

// GetAllTrades retrieves all trades
func GetAllTrades() ([]*models.Trade, error) {
	var trades []*models.Trade
	err := database.GetByPrefix(TRADE_PREFIX, &trades)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	return trades, nil
}

// DeleteTrade deletes a trade by ID
func DeleteTrade(id string) error {
	return database.Delete(id)
}
