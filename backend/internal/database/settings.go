package database

import (
	"encoding/json"
	"time"

	"github.com/spotify-genre-organizer/backend/internal/models"
)

// GetUserSettings fetches settings for a user, returning defaults if not found
func GetUserSettings(userID string) (*models.UserSettings, error) {
	// Try to fetch from DB
	res, _, err := Client.From("user_settings").
		Select("*", "", false).
		Eq("user_id", userID).
		Single().
		Execute()

	if err != nil {
		// If error is "no rows" or similar, return default settings
		return models.DefaultSettings(userID), nil
	}

	var settings models.UserSettings
	if err := json.Unmarshal(res, &settings); err != nil {
		// If unmarshal fails, maybe empty response? Return default
		return models.DefaultSettings(userID), nil
	}

	return &settings, nil
}

// SaveUserSettings upserts user settings
func SaveUserSettings(settings *models.UserSettings) error {
	settings.UpdatedAt = time.Now()

	_, _, err := Client.From("user_settings").
		Upsert(settings, "", "", "").
		Execute()

	return err
}

// GetPlaylistOverrides fetches all overrides for a user
func GetPlaylistOverrides(userID string) (map[string]*models.PlaylistOverride, error) {
	res, _, err := Client.From("playlist_overrides").
		Select("*", "", false).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return nil, err
	}

	var overrides []models.PlaylistOverride
	if err := json.Unmarshal(res, &overrides); err != nil {
		return nil, err
	}

	result := make(map[string]*models.PlaylistOverride)
	for i := range overrides {
		result[overrides[i].PlaylistSpotifyID] = &overrides[i]
	}

	return result, nil
}

// SavePlaylistOverride upserts a playlist override
func SavePlaylistOverride(override *models.PlaylistOverride) error {
	override.UpdatedAt = time.Now()

	_, _, err := Client.From("playlist_overrides").
		Upsert(override, "", "", "").
		Execute()

	return err
}
