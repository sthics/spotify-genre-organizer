package database

import (
	"encoding/json"
	"time"

	"github.com/spotify-genre-organizer/backend/internal/models"
)

// GetExclusionRules returns all exclusion rules for a user
func GetExclusionRules(userID string) ([]models.ExclusionRule, error) {
	res, _, err := Client.From("exclusion_rules").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", nil).
		Execute()

	if err != nil {
		return nil, err
	}

	var rules []models.ExclusionRule
	if err := json.Unmarshal(res, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateExclusionRule creates a new exclusion rule
func CreateExclusionRule(rule *models.ExclusionRule) error {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	res, _, err := Client.From("exclusion_rules").
		Insert(rule, false, "", "", "").
		Execute()

	if err != nil {
		return err
	}

	// Parse response to get the generated ID
	var created []models.ExclusionRule
	if err := json.Unmarshal(res, &created); err != nil {
		return err
	}
	if len(created) > 0 {
		rule.ID = created[0].ID
	}

	return nil
}

// DeleteExclusionRule removes an exclusion rule by ID
func DeleteExclusionRule(userID, ruleID string) error {
	_, _, err := Client.From("exclusion_rules").
		Delete("", "").
		Eq("id", ruleID).
		Eq("user_id", userID).
		Execute()

	return err
}

// GetExclusionRuleByID fetches a single rule
func GetExclusionRuleByID(userID, ruleID string) (*models.ExclusionRule, error) {
	res, _, err := Client.From("exclusion_rules").
		Select("*", "", false).
		Eq("id", ruleID).
		Eq("user_id", userID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var rule models.ExclusionRule
	if err := json.Unmarshal(res, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}
