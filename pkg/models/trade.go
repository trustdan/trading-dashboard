package models

import "time"

// Trade represents a trading position
type Trade struct {
	ID             int       `json:"id"`
	EntryDate      time.Time `json:"entryDate"`
	Ticker         string    `json:"ticker"`
	Sector         string    `json:"sector"`
	EntryPrice     float64   `json:"entryPrice"`
	Notes          string    `json:"notes"`
	ExpirationDate time.Time `json:"expirationDate"`
	StrategyType   string    `json:"strategyType"`
	SpreadType     string    `json:"spreadType"`
	Direction      string    `json:"direction"`
}
