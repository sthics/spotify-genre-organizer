package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Song struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
	Genres  []string `json:"genres"`
}

type Artist struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type likedSongsResponse struct {
	Items []struct {
		Track struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"track"`
	} `json:"items"`
	Total int     `json:"total"`
	Next  *string `json:"next"`
}

func ParseLikedSongsResponse(data []byte) ([]Song, int, string, error) {
	var resp likedSongsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, 0, "", err
	}

	songs := make([]Song, len(resp.Items))
	for i, item := range resp.Items {
		artists := make([]Artist, len(item.Track.Artists))
		for j, a := range item.Track.Artists {
			artists[j] = Artist{
				ID:   a.ID,
				Name: a.Name,
			}
		}
		songs[i] = Song{
			ID:      item.Track.ID,
			Name:    item.Track.Name,
			Artists: artists,
		}
	}

	next := ""
	if resp.Next != nil {
		next = *resp.Next
	}

	return songs, resp.Total, next, nil
}

func FetchLikedSongs(accessToken string, limit, offset int) ([]Song, int, string, error) {
	url := fmt.Sprintf("%s/me/tracks?limit=%d&offset=%d", APIURL, limit, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, "", fmt.Errorf("failed to fetch liked songs: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}

	return ParseLikedSongsResponse(body)
}

func FetchAllLikedSongs(accessToken string, progressCallback func(processed, total int)) ([]Song, error) {
	var allSongs []Song
	limit := 50
	offset := 0
	total := 0

	for {
		songs, t, _, err := FetchLikedSongs(accessToken, limit, offset)
		if err != nil {
			return nil, err
		}

		if total == 0 {
			total = t
		}

		allSongs = append(allSongs, songs...)

		if progressCallback != nil {
			progressCallback(len(allSongs), total)
		}

		if len(songs) < limit || len(allSongs) >= total {
			break
		}

		offset += limit
	}

	return allSongs, nil
}

func EnrichSongsWithGenres(songs []Song, artistGenres map[string][]string) {
	for i := range songs {
		genreSet := make(map[string]bool)
		for _, artist := range songs[i].Artists {
			if genres, ok := artistGenres[artist.ID]; ok {
				for _, g := range genres {
					genreSet[g] = true
				}
			}
		}

		songs[i].Genres = make([]string, 0, len(genreSet))
		for g := range genreSet {
			songs[i].Genres = append(songs[i].Genres, g)
		}
	}
}
