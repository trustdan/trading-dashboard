package repositories

import (
	"fmt"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

// SaveRiskAssessment saves a risk assessment to the database
func SaveRiskAssessment(assessment *models.RiskAssessment) error {
	// Calculate the overall score
	assessment.CalculateOverallScore()

	var query string
	var args []interface{}

	if assessment.ID == 0 {
		// Insert new record
		query = `
		INSERT INTO risk_assessments 
		(date, emotional, fomo, bias, physical, pnl, overall_score) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`
		args = []interface{}{
			assessment.Date.Format("2006-01-02"),
			assessment.Emotional,
			assessment.Fomo,
			assessment.Bias,
			assessment.Physical,
			assessment.Pnl,
			assessment.OverallScore,
		}
	} else {
		// Update existing record
		query = `
		UPDATE risk_assessments 
		SET date = ?, emotional = ?, fomo = ?, bias = ?, physical = ?, pnl = ?, overall_score = ? 
		WHERE id = ?`
		args = []interface{}{
			assessment.Date.Format("2006-01-02"),
			assessment.Emotional,
			assessment.Fomo,
			assessment.Bias,
			assessment.Physical,
			assessment.Pnl,
			assessment.OverallScore,
			assessment.ID,
		}
	}

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to save risk assessment: %w", err)
	}

	// If it was an insert, get the ID
	if assessment.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
		assessment.ID = int(id)
	}

	return nil
}

// GetRiskAssessment retrieves a risk assessment by ID
func GetRiskAssessment(id int) (*models.RiskAssessment, error) {
	query := "SELECT id, date, emotional, fomo, bias, physical, pnl, overall_score FROM risk_assessments WHERE id = ?"

	var dateStr string
	assessment := &models.RiskAssessment{}

	err := database.DB.QueryRow(query, id).Scan(
		&assessment.ID,
		&dateStr,
		&assessment.Emotional,
		&assessment.Fomo,
		&assessment.Bias,
		&assessment.Physical,
		&assessment.Pnl,
		&assessment.OverallScore,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get risk assessment: %w", err)
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	assessment.Date = date

	return assessment, nil
}

// GetLatestRiskAssessment retrieves the latest risk assessment
func GetLatestRiskAssessment() (*models.RiskAssessment, error) {
	query := "SELECT id, date, emotional, fomo, bias, physical, pnl, overall_score FROM risk_assessments ORDER BY date DESC LIMIT 1"

	var dateStr string
	assessment := &models.RiskAssessment{}

	err := database.DB.QueryRow(query).Scan(
		&assessment.ID,
		&dateStr,
		&assessment.Emotional,
		&assessment.Fomo,
		&assessment.Bias,
		&assessment.Physical,
		&assessment.Pnl,
		&assessment.OverallScore,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get latest risk assessment: %w", err)
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	assessment.Date = date

	return assessment, nil
}

// GetAllRiskAssessments retrieves all risk assessments
func GetAllRiskAssessments() ([]*models.RiskAssessment, error) {
	query := "SELECT id, date, emotional, fomo, bias, physical, pnl, overall_score FROM risk_assessments ORDER BY date DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query risk assessments: %w", err)
	}
	defer rows.Close()

	var assessments []*models.RiskAssessment

	for rows.Next() {
		var dateStr string
		assessment := &models.RiskAssessment{}

		err := rows.Scan(
			&assessment.ID,
			&dateStr,
			&assessment.Emotional,
			&assessment.Fomo,
			&assessment.Bias,
			&assessment.Physical,
			&assessment.Pnl,
			&assessment.OverallScore,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan risk assessment row: %w", err)
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		assessment.Date = date

		assessments = append(assessments, assessment)
	}

	return assessments, nil
}
