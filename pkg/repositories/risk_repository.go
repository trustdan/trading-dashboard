package repositories

import (
	"fmt"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
)

const RISK_PREFIX = "risk_"

// SaveRiskAssessment saves a risk assessment to the database
func SaveRiskAssessment(assessment *models.RiskAssessment) error {
	// Calculate the overall score
	assessment.CalculateOverallScore()

	// If no ID is set, generate one
	if assessment.ID == "" {
		assessment.ID = database.GenerateKey(RISK_PREFIX)
	}

	// Save the assessment to BadgerDB
	return database.Set(assessment.ID, assessment)
}

// GetRiskAssessment retrieves a risk assessment by ID
func GetRiskAssessment(id string) (*models.RiskAssessment, error) {
	assessment := &models.RiskAssessment{}
	err := database.Get(id, assessment)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk assessment: %w", err)
	}
	return assessment, nil
}

// GetLatestRiskAssessment retrieves the latest risk assessment
func GetLatestRiskAssessment() (*models.RiskAssessment, error) {
	assessments, err := GetAllRiskAssessments()
	if err != nil {
		return nil, err
	}

	if len(assessments) == 0 {
		return nil, fmt.Errorf("no risk assessments found")
	}

	// Since we're getting all assessments, sort them by date (newest first)
	// The first one will be the latest
	latestAssessment := assessments[0]
	latestDate := latestAssessment.Date

	for _, assessment := range assessments {
		if assessment.Date.After(latestDate) {
			latestAssessment = assessment
			latestDate = assessment.Date
		}
	}

	return latestAssessment, nil
}

// GetAllRiskAssessments retrieves all risk assessments
func GetAllRiskAssessments() ([]*models.RiskAssessment, error) {
	var assessments []*models.RiskAssessment
	err := database.GetByPrefix(RISK_PREFIX, &assessments)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk assessments: %w", err)
	}
	return assessments, nil
}

// DeleteRiskAssessment deletes a risk assessment by ID
func DeleteRiskAssessment(id string) error {
	return database.Delete(id)
}
