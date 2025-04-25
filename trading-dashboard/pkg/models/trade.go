package models

import "time"

// Trade represents a trading position entry
type Trade struct {
	ID         int       `json:"id"`
	EntryDate  time.Time `json:"entryDate"`
	Ticker     string    `json:"ticker"`
	Sector     string    `json:"sector"`
	EntryPrice float64   `json:"entryPrice"`
	Notes      string    `json:"notes"`
}
