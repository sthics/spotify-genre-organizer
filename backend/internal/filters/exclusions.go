package filters

import (
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

// ApplyExclusions filters out songs matching exclusion rules
func ApplyExclusions(songs []spotify.Song, rules []models.ExclusionRule) []spotify.Song {
	if len(rules) == 0 {
		return songs
	}

	// Build lookup maps for O(1) checks
	blockedSongs := make(map[string]bool)
	blockedArtistsAll := make(map[string]bool)
	blockedArtistsPrimary := make(map[string]bool)

	for _, rule := range rules {
		switch rule.ExclusionType {
		case models.ExclusionTypeSong:
			blockedSongs[rule.SpotifyID] = true
		case models.ExclusionTypeArtist:
			if rule.Scope != nil && *rule.Scope == models.ExclusionScopeAllAppearances {
				blockedArtistsAll[rule.SpotifyID] = true
			} else {
				blockedArtistsPrimary[rule.SpotifyID] = true
			}
		}
	}

	// Filter songs
	result := make([]spotify.Song, 0, len(songs))
	for _, song := range songs {
		if blockedSongs[song.ID] {
			continue
		}

		// Check primary artist (first in list)
		if len(song.Artists) > 0 {
			primaryID := song.Artists[0].ID
			if blockedArtistsPrimary[primaryID] || blockedArtistsAll[primaryID] {
				continue
			}
		}

		// Check all artists for "all_appearances" scope
		blocked := false
		for _, artist := range song.Artists {
			if blockedArtistsAll[artist.ID] {
				blocked = true
				break
			}
		}
		if blocked {
			continue
		}

		result = append(result, song)
	}

	return result
}
