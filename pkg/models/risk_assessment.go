package models

import "time"

// RiskAssessment represents a daily risk assessment entry
type RiskAssessment struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	Emotional    int       `json:"emotional"`    // Range: -3 to +3
	Fomo         int       `json:"fomo"`         // Range: -3 to +3
	Bias         int       `json:"bias"`         // Range: -3 to +3
	Physical     int       `json:"physical"`     // Range: -3 to +3
	Pnl          int       `json:"pnl"`          // Range: -3 to +3
	OverallScore int       `json:"overallScore"` // Calculated score
}

// CalculateOverallScore calculates the overall risk score
func (ra *RiskAssessment) CalculateOverallScore() {
	// Simple average for now, can be enhanced with weighted calculations
	ra.OverallScore = (ra.Emotional + ra.Fomo + ra.Bias + ra.Physical + ra.Pnl) / 5
}
