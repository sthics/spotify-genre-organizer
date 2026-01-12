package spotify

import (
	"testing"
)

func TestBuildPlaylistName(t *testing.T) {
	tests := []struct {
		genre    string
		expected string
	}{
		{"Rock", "Rock by Organizer"},
		{"Hip-Hop", "Hip-Hop by Organizer"},
		{"R&B", "R&B by Organizer"},
	}

	for _, tt := range tests {
		result := BuildPlaylistName(tt.genre)
		if result != tt.expected {
			t.Errorf("BuildPlaylistName(%q) = %q, want %q", tt.genre, result, tt.expected)
		}
	}
}

func TestChunkTrackIDs(t *testing.T) {
	ids := make([]string, 150)
	for i := range ids {
		ids[i] = "track" + string(rune('0'+i%10))
	}

	chunks := ChunkTrackIDs(ids, 100)

	if len(chunks) != 2 {
		t.Errorf("expected 2 chunks, got %d", len(chunks))
	}

	if len(chunks[0]) != 100 {
		t.Errorf("expected first chunk to have 100 items, got %d", len(chunks[0]))
	}

	if len(chunks[1]) != 50 {
		t.Errorf("expected second chunk to have 50 items, got %d", len(chunks[1]))
	}
}
