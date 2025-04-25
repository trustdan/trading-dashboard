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
		(entry_date, ticker, sector, entry_price, notes, expiration_date, strategy_type, spread_type, direction) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		args = []interface{}{
			trade.EntryDate.Format("2006-01-02"),
			trade.Ticker,
			trade.Sector,
			trade.EntryPrice,
			trade.Notes,
			trade.ExpirationDate.Format("2006-01-02"),
			trade.StrategyType,
			trade.SpreadType,
			trade.Direction,
		}
	} else {
		// Update existing record
		query = `
		UPDATE trades 
		SET entry_date = ?, ticker = ?, sector = ?, entry_price = ?, notes = ?,
		    expiration_date = ?, strategy_type = ?, spread_type = ?, direction = ?
		WHERE id = ?`
		args = []interface{}{
			trade.EntryDate.Format("2006-01-02"),
			trade.Ticker,
			trade.Sector,
			trade.EntryPrice,
			trade.Notes,
			trade.ExpirationDate.Format("2006-01-02"),
			trade.StrategyType,
			trade.SpreadType,
			trade.Direction,
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
	query := `
	SELECT id, entry_date, ticker, sector, entry_price, notes, 
	       expiration_date, strategy_type, spread_type, direction 
	FROM trades 
	WHERE id = ?`

	var dateStr, expirationStr string
	trade := &models.Trade{}

	err := database.DB.QueryRow(query, id).Scan(
		&trade.ID,
		&dateStr,
		&trade.Ticker,
		&trade.Sector,
		&trade.EntryPrice,
		&trade.Notes,
		&expirationStr,
		&trade.StrategyType,
		&trade.SpreadType,
		&trade.Direction,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	// Parse the dates
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry date: %w", err)
	}
	trade.EntryDate = date

	// Parse expiration date if it exists
	if expirationStr != "" {
		expiration, err := time.Parse("2006-01-02", expirationStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration date: %w", err)
		}
		trade.ExpirationDate = expiration
	}

	return trade, nil
}

// GetTradesByDateRange retrieves trades within a date range
func GetTradesByDateRange(startDate, endDate time.Time) ([]*models.Trade, error) {
	query := `
	SELECT id, entry_date, ticker, sector, entry_price, notes, 
	       expiration_date, strategy_type, spread_type, direction 
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
		var dateStr, expirationStr string
		trade := &models.Trade{}

		err := rows.Scan(
			&trade.ID,
			&dateStr,
			&trade.Ticker,
			&trade.Sector,
			&trade.EntryPrice,
			&trade.Notes,
			&expirationStr,
			&trade.StrategyType,
			&trade.SpreadType,
			&trade.Direction,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		// Parse the dates
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse entry date: %w", err)
		}
		trade.EntryDate = date

		// Parse expiration date if it exists
		if expirationStr != "" {
			expiration, err := time.Parse("2006-01-02", expirationStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expiration date: %w", err)
			}
			trade.ExpirationDate = expiration
		}

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
	query := `
	SELECT id, entry_date, ticker, sector, entry_price, notes, 
	       expiration_date, strategy_type, spread_type, direction 
	FROM trades 
	ORDER BY entry_date DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*models.Trade

	for rows.Next() {
		var dateStr, expirationStr string
		trade := &models.Trade{}

		err := rows.Scan(
			&trade.ID,
			&dateStr,
			&trade.Ticker,
			&trade.Sector,
			&trade.EntryPrice,
			&trade.Notes,
			&expirationStr,
			&trade.StrategyType,
			&trade.SpreadType,
			&trade.Direction,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		// Parse the dates
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse entry date: %w", err)
		}
		trade.EntryDate = date

		// Parse expiration date if it exists
		if expirationStr != "" {
			expiration, err := time.Parse("2006-01-02", expirationStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expiration date: %w", err)
			}
			trade.ExpirationDate = expiration
		}

		trades = append(trades, trade)
	}

	return trades, nil
}

// DeleteTrade removes a trade from the database
func DeleteTrade(id int) error {
	query := "DELETE FROM trades WHERE id = ?"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no trade found with ID %d to delete", id)
	}

	return nil
}
