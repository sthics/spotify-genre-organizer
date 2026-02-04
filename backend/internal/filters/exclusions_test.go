package filters

import (
	"testing"

	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

func TestApplyExclusions_BlocksArtistAllAppearances(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}}},
		{ID: "3", Name: "Collab", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}, {ID: "artist3", Name: "Artist Three"}}},
	}

	allAppearances := models.ExclusionScopeAllAppearances
	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeArtist, SpotifyID: "artist2", Scope: &allAppearances},
	}

	result := ApplyExclusions(songs, rules)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("expected song 1, got %s", result[0].ID)
	}
}

func TestApplyExclusions_BlocksArtistPrimaryOnly(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}}},
		{ID: "3", Name: "Collab", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}, {ID: "artist2", Name: "Artist Two"}}},
	}

	primaryOnly := models.ExclusionScopePrimaryOnly
	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeArtist, SpotifyID: "artist2", Scope: &primaryOnly},
	}

	result := ApplyExclusions(songs, rules)

	// Song 2 blocked (artist2 is primary), Song 3 kept (artist2 is feature)
	if len(result) != 2 {
		t.Errorf("expected 2 songs, got %d", len(result))
	}
}

func TestApplyExclusions_BlocksSong(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
	}

	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeSong, SpotifyID: "1"},
	}

	result := ApplyExclusions(songs, rules)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
	if result[0].ID != "2" {
		t.Errorf("expected song 2, got %s", result[0].ID)
	}
}

func TestApplyExclusions_NoRules(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
	}

	result := ApplyExclusions(songs, nil)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
}
