package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ArtistDetails struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type artistsResponse struct {
	Artists []ArtistDetails `json:"artists"`
}

func FetchArtists(accessToken string, artistIDs []string) ([]ArtistDetails, error) {
	if len(artistIDs) == 0 {
		return nil, nil
	}

	if len(artistIDs) > 50 {
		artistIDs = artistIDs[:50]
	}

	url := fmt.Sprintf("%s/artists?ids=%s", APIURL, strings.Join(artistIDs, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch artists: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result artistsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Artists, nil
}

func FetchAllArtistGenres(accessToken string, songs []Song, progressCallback func(processed, total int)) (map[string][]string, error) {
	artistSet := make(map[string]bool)
	for _, song := range songs {
		for _, artist := range song.Artists {
			artistSet[artist.ID] = true
		}
	}

	artistIDs := make([]string, 0, len(artistSet))
	for id := range artistSet {
		artistIDs = append(artistIDs, id)
	}

	genreMap := make(map[string][]string)
	batchSize := 50
	total := len(artistIDs)

	for i := 0; i < len(artistIDs); i += batchSize {
		end := i + batchSize
		if end > len(artistIDs) {
			end = len(artistIDs)
		}

		batch := artistIDs[i:end]
		artists, err := FetchArtists(accessToken, batch)
		if err != nil {
			return nil, err
		}

		for _, artist := range artists {
			genreMap[artist.ID] = artist.Genres
		}

		if progressCallback != nil {
			progressCallback(end, total)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return genreMap, nil
}
