package models

import "time"

// StockRating represents a rating for a stock
type StockRating struct {
	ID                    int       `json:"id"`
	Date                  time.Time `json:"date"`
	Ticker                string    `json:"ticker"`
	MarketSentiment       int       `json:"marketSentiment"`       // Range: -3 to +3
	BasicMaterials        int       `json:"basicMaterials"`        // Range: -3 to +3
	CommunicationServices int       `json:"communicationServices"` // Range: -3 to +3
	ConsumerCyclical      int       `json:"consumerCyclical"`      // Range: -3 to +3
	ConsumerDefensive     int       `json:"consumerDefensive"`     // Range: -3 to +3
	Energy                int       `json:"energy"`                // Range: -3 to +3
	Financial             int       `json:"financial"`             // Range: -3 to +3
	Healthcare            int       `json:"healthcare"`            // Range: -3 to +3
	Industrials           int       `json:"industrials"`           // Range: -3 to +3
	RealEstate            int       `json:"realEstate"`            // Range: -3 to +3
	Technology            int       `json:"technology"`            // Range: -3 to +3
	Utilities             int       `json:"utilities"`             // Range: -3 to +3
	StockSentiment        int       `json:"stockSentiment"`        // Range: -3 to +3
	Pattern               string    `json:"pattern"`               // Chart pattern
	EnthusiasmRating      int       `json:"enthusiasmRating"`      // Calculated rating
}

// CalculateEnthusiasm calculates the enthusiasm rating
func (sr *StockRating) CalculateEnthusiasm() {
	// Map of pattern values
	patternPoints := map[string]int{
		// Original patterns
		"High Base":           2,
		"Low Base":            2,
		"Ascending Triangle":  3,
		"Descending Triangle": 3,
		"Bull Pullback":       2,
		"Bear Rally":          2,
		"Double-Top":          3,
		"Cup-and-Handle":      4,
		// Additional patterns
		"Head and Shoulders":         4,
		"Inverse Head and Shoulders": 4,
		"Bullish Flag":               3,
		"Bearish Flag":               3,
		"Rising Wedge":               2,
		"Falling Wedge":              2,
		"Double Bottom":              3,
		"Rounding Bottom":            3,
		"Breakaway Gap":              3,
		"Runaway Gap":                2,
		"Exhaustion Gap":             1,
		"Bullish Engulfing":          2,
		"Bearish Engulfing":          2,
	}

	// Get pattern points, default to 0 if pattern not found
	points, exists := patternPoints[sr.Pattern]
	if !exists {
		points = 0
	}

	// Calculate enthusiasm based on stock sentiment and pattern
	sr.EnthusiasmRating = sr.StockSentiment + points
}
