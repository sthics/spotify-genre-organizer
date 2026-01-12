package spotify

import (
	"testing"
)

func TestParseLikedSongsResponse(t *testing.T) {
	jsonData := `{
		"items": [
			{
				"track": {
					"id": "track123",
					"name": "Test Song",
					"artists": [
						{
							"id": "artist123",
							"name": "Test Artist"
						}
					]
				}
			}
		],
		"total": 1,
		"next": null
	}`

	songs, total, next, err := ParseLikedSongsResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}

	if next != "" {
		t.Errorf("expected empty next, got %s", next)
	}

	if len(songs) != 1 {
		t.Fatalf("expected 1 song, got %d", len(songs))
	}

	if songs[0].ID != "track123" {
		t.Errorf("expected track ID track123, got %s", songs[0].ID)
	}
}
