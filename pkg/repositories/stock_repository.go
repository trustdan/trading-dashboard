package repositories

import (
	"fmt"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

// SaveStockRating saves a stock rating to the database
func SaveStockRating(rating *models.StockRating) error {
	// Calculate the enthusiasm rating
	rating.CalculateEnthusiasm()

	var query string
	var args []interface{}

	if rating.ID == 0 {
		// Insert new record
		query = `
		INSERT INTO stock_ratings 
		(date, ticker, market_sentiment, basic_materials, communication_services, 
		 consumer_cyclical, consumer_defensive, energy, financial, healthcare, 
		 industrials, real_estate, technology, utilities, stock_sentiment, pattern, enthusiasm_rating) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		args = []interface{}{
			rating.Date.Format("2006-01-02"),
			rating.Ticker,
			rating.MarketSentiment,
			rating.BasicMaterials,
			rating.CommunicationServices,
			rating.ConsumerCyclical,
			rating.ConsumerDefensive,
			rating.Energy,
			rating.Financial,
			rating.Healthcare,
			rating.Industrials,
			rating.RealEstate,
			rating.Technology,
			rating.Utilities,
			rating.StockSentiment,
			rating.Pattern,
			rating.EnthusiasmRating,
		}
	} else {
		// Update existing record
		query = `
		UPDATE stock_ratings 
		SET date = ?, ticker = ?, market_sentiment = ?, 
		    basic_materials = ?, communication_services = ?, 
		    consumer_cyclical = ?, consumer_defensive = ?, 
		    energy = ?, financial = ?, healthcare = ?, 
		    industrials = ?, real_estate = ?, technology = ?, utilities = ?,
		    stock_sentiment = ?, pattern = ?, enthusiasm_rating = ? 
		WHERE id = ?`
		args = []interface{}{
			rating.Date.Format("2006-01-02"),
			rating.Ticker,
			rating.MarketSentiment,
			rating.BasicMaterials,
			rating.CommunicationServices,
			rating.ConsumerCyclical,
			rating.ConsumerDefensive,
			rating.Energy,
			rating.Financial,
			rating.Healthcare,
			rating.Industrials,
			rating.RealEstate,
			rating.Technology,
			rating.Utilities,
			rating.StockSentiment,
			rating.Pattern,
			rating.EnthusiasmRating,
			rating.ID,
		}
	}

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to save stock rating: %w", err)
	}

	// If it was an insert, get the ID
	if rating.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
		rating.ID = int(id)
	}

	return nil
}

// GetStockRating retrieves a stock rating by ID
func GetStockRating(id int) (*models.StockRating, error) {
	query := `
	SELECT id, date, ticker, market_sentiment, 
	       basic_materials, communication_services, 
	       consumer_cyclical, consumer_defensive, 
	       energy, financial, healthcare, 
	       industrials, real_estate, technology, utilities,
	       stock_sentiment, pattern, enthusiasm_rating 
	FROM stock_ratings WHERE id = ?`

	var dateStr string
	rating := &models.StockRating{}

	err := database.DB.QueryRow(query, id).Scan(
		&rating.ID,
		&dateStr,
		&rating.Ticker,
		&rating.MarketSentiment,
		&rating.BasicMaterials,
		&rating.CommunicationServices,
		&rating.ConsumerCyclical,
		&rating.ConsumerDefensive,
		&rating.Energy,
		&rating.Financial,
		&rating.Healthcare,
		&rating.Industrials,
		&rating.RealEstate,
		&rating.Technology,
		&rating.Utilities,
		&rating.StockSentiment,
		&rating.Pattern,
		&rating.EnthusiasmRating,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stock rating: %w", err)
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	rating.Date = date

	return rating, nil
}

// GetStockRatingsByTicker retrieves all stock ratings for a specific ticker
func GetStockRatingsByTicker(ticker string) ([]*models.StockRating, error) {
	query := `
	SELECT id, date, ticker, market_sentiment, 
	       basic_materials, communication_services, 
	       consumer_cyclical, consumer_defensive, 
	       energy, financial, healthcare, 
	       industrials, real_estate, technology, utilities,
	       stock_sentiment, pattern, enthusiasm_rating 
	FROM stock_ratings 
	WHERE ticker = ? 
	ORDER BY date DESC`

	rows, err := database.DB.Query(query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock ratings: %w", err)
	}
	defer rows.Close()

	var ratings []*models.StockRating

	for rows.Next() {
		var dateStr string
		rating := &models.StockRating{}

		err := rows.Scan(
			&rating.ID,
			&dateStr,
			&rating.Ticker,
			&rating.MarketSentiment,
			&rating.BasicMaterials,
			&rating.CommunicationServices,
			&rating.ConsumerCyclical,
			&rating.ConsumerDefensive,
			&rating.Energy,
			&rating.Financial,
			&rating.Healthcare,
			&rating.Industrials,
			&rating.RealEstate,
			&rating.Technology,
			&rating.Utilities,
			&rating.StockSentiment,
			&rating.Pattern,
			&rating.EnthusiasmRating,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan stock rating row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		rating.Date = date

		ratings = append(ratings, rating)
	}

	return ratings, nil
}

// GetAllStockRatings retrieves all stock ratings
func GetAllStockRatings() ([]*models.StockRating, error) {
	query := `
	SELECT id, date, ticker, market_sentiment, 
	       basic_materials, communication_services, 
	       consumer_cyclical, consumer_defensive, 
	       energy, financial, healthcare, 
	       industrials, real_estate, technology, utilities,
	       stock_sentiment, pattern, enthusiasm_rating 
	FROM stock_ratings 
	ORDER BY date DESC, enthusiasm_rating DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock ratings: %w", err)
	}
	defer rows.Close()

	var ratings []*models.StockRating

	for rows.Next() {
		var dateStr string
		rating := &models.StockRating{}

		err := rows.Scan(
			&rating.ID,
			&dateStr,
			&rating.Ticker,
			&rating.MarketSentiment,
			&rating.BasicMaterials,
			&rating.CommunicationServices,
			&rating.ConsumerCyclical,
			&rating.ConsumerDefensive,
			&rating.Energy,
			&rating.Financial,
			&rating.Healthcare,
			&rating.Industrials,
			&rating.RealEstate,
			&rating.Technology,
			&rating.Utilities,
			&rating.StockSentiment,
			&rating.Pattern,
			&rating.EnthusiasmRating,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan stock rating row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		rating.Date = date

		ratings = append(ratings, rating)
	}

	return ratings, nil
}
