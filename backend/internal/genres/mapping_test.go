package genres

import (
	"testing"
)

func TestConsolidateGenre(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"indie rock", "Rock"},
		{"alternative rock", "Rock"},
		{"classic rock", "Rock"},
		{"edm", "Electronic"},
		{"house", "Electronic"},
		{"hip hop", "Hip-Hop"},
		{"rap", "Hip-Hop"},
		{"jazz fusion", "Jazz"},
		{"unknown genre xyz", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ConsolidateGenre(tt.input)
			if result != tt.expected {
				t.Errorf("ConsolidateGenre(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetParentGenres(t *testing.T) {
	genres := GetParentGenres()

	if len(genres) < 10 {
		t.Errorf("expected at least 10 parent genres, got %d", len(genres))
	}

	expected := []string{"Rock", "Pop", "Hip-Hop", "Electronic", "Jazz", "Classical"}
	for _, e := range expected {
		found := false
		for _, g := range genres {
			if g == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected parent genre %q not found", e)
		}
	}
}

func TestGenrePriority(t *testing.T) {
	// Reggae should come before Hip-Hop in priority
	reggaeIdx := -1
	hipHopIdx := -1
	for i, g := range GenrePriority {
		if g == "Reggae" {
			reggaeIdx = i
		}
		if g == "Hip-Hop" {
			hipHopIdx = i
		}
	}
	if reggaeIdx == -1 || hipHopIdx == -1 {
		t.Fatal("GenrePriority missing Reggae or Hip-Hop")
	}
	if reggaeIdx >= hipHopIdx {
		t.Errorf("Reggae (idx %d) should come before Hip-Hop (idx %d)", reggaeIdx, hipHopIdx)
	}
}

func TestScoreGenres(t *testing.T) {
	tests := []struct {
		name     string
		genres   []string
		expected string
	}{
		{
			name:     "single genre",
			genres:   []string{"rock"},
			expected: "Rock",
		},
		{
			name:     "dub should win over hip hop",
			genres:   []string{"hip hop", "dub", "reggae"},
			expected: "Reggae", // 2 votes vs 1
		},
		{
			name:     "electronic wins with multiple matches",
			genres:   []string{"dub", "trip hop", "downtempo", "acid jazz"},
			expected: "Electronic", // trip hop + downtempo = 2 votes
		},
		{
			name:     "tie broken by priority - reggae before hip-hop",
			genres:   []string{"hip hop", "dub"},
			expected: "Reggae", // 1-1 tie, reggae has priority
		},
		{
			name:     "empty genres",
			genres:   []string{},
			expected: "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScoreGenres(tt.genres)
			if result != tt.expected {
				t.Errorf("ScoreGenres(%v) = %q, want %q", tt.genres, result, tt.expected)
			}
		})
	}
}
