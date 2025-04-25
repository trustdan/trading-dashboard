package repositories

import (
	"fmt"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

// SaveTrade saves a trade to the database
func SaveTrade(trade *models.Trade) error {
	var query string
	var args []interface{}

	if trade.ID == 0 {
		// Insert new record
		query = `
		INSERT INTO trades 
		(entry_date, ticker, sector, entry_price, notes) 
		VALUES (?, ?, ?, ?, ?)`
		args = []interface{}{
			trade.EntryDate.Format("2006-01-02"),
			trade.Ticker,
			trade.Sector,
			trade.EntryPrice,
			trade.Notes,
		}
	} else {
		// Update existing record
		query = `
		UPDATE trades 
		SET entry_date = ?, ticker = ?, sector = ?, entry_price = ?, notes = ? 
		WHERE id = ?`
		args = []interface{}{
			trade.EntryDate.Format("2006-01-02"),
			trade.Ticker,
			trade.Sector,
			trade.EntryPrice,
			trade.Notes,
			trade.ID,
		}
	}

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to save trade: %w", err)
	}

	// If it was an insert, get the ID
	if trade.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
		trade.ID = int(id)
	}

	return nil
}

// GetTrade retrieves a trade by ID
func GetTrade(id int) (*models.Trade, error) {
	query := "SELECT id, entry_date, ticker, sector, entry_price, notes FROM trades WHERE id = ?"

	var dateStr string
	trade := &models.Trade{}

	err := database.DB.QueryRow(query, id).Scan(
		&trade.ID,
		&dateStr,
		&trade.Ticker,
		&trade.Sector,
		&trade.EntryPrice,
		&trade.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	trade.EntryDate = date

	return trade, nil
}

// GetTradesByDateRange retrieves trades within a date range
func GetTradesByDateRange(startDate, endDate time.Time) ([]*models.Trade, error) {
	query := `
	SELECT id, entry_date, ticker, sector, entry_price, notes 
	FROM trades 
	WHERE entry_date BETWEEN ? AND ? 
	ORDER BY entry_date DESC`

	rows, err := database.DB.Query(query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*models.Trade

	for rows.Next() {
		var dateStr string
		trade := &models.Trade{}

		err := rows.Scan(
			&trade.ID,
			&dateStr,
			&trade.Ticker,
			&trade.Sector,
			&trade.EntryPrice,
			&trade.Notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		trade.EntryDate = date

		trades = append(trades, trade)
	}

	return trades, nil
}

// GetTradesByTicker retrieves trades for a specific ticker
func GetTradesByTicker(ticker string) ([]*models.Trade, error) {
	query := `
	SELECT id, entry_date, ticker, sector, entry_price, notes 
	FROM trades 
	WHERE ticker = ? 
	ORDER BY entry_date DESC`

	rows, err := database.DB.Query(query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*models.Trade

	for rows.Next() {
		var dateStr string
		trade := &models.Trade{}

		err := rows.Scan(
			&trade.ID,
			&dateStr,
			&trade.Ticker,
			&trade.Sector,
			&trade.EntryPrice,
			&trade.Notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		trade.EntryDate = date

		trades = append(trades, trade)
	}

	return trades, nil
}

// GetAllTrades retrieves all trades
func GetAllTrades() ([]*models.Trade, error) {
	query := "SELECT id, entry_date, ticker, sector, entry_price, notes FROM trades ORDER BY entry_date DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*models.Trade

	for rows.Next() {
		var dateStr string
		trade := &models.Trade{}

		err := rows.Scan(
			&trade.ID,
			&dateStr,
			&trade.Ticker,
			&trade.Sector,
			&trade.EntryPrice,
			&trade.Notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		trade.EntryDate = date

		trades = append(trades, trade)
	}

	return trades, nil
}
