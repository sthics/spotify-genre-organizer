package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Playlist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_urls"`
	TracksTotal int    `json:"tracks_total"`
}

type PlaylistsResponse struct {
	Items []PlaylistItem `json:"items"`
	Total int            `json:"total"`
}

type PlaylistItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Tracks struct {
		Total int `json:"total"`
	} `json:"tracks"`
	Owner struct {
		ID string `json:"id"`
	} `json:"owner"`
}

type createPlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type addTracksRequest struct {
	URIs []string `json:"uris"`
}

func BuildPlaylistName(genre string) string {
	return genre + " by Organizer"
}

func ChunkTrackIDs(ids []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunks = append(chunks, ids[i:end])
	}
	return chunks
}

func CreatePlaylist(accessToken, userID, name, description string) (*Playlist, error) {
	url := fmt.Sprintf("%s/users/%s/playlists", APIURL, userID)

	body := createPlaylistRequest{
		Name:        name,
		Description: description,
		Public:      false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create playlist: %d - %s", resp.StatusCode, string(respBody))
	}

	var playlist struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return nil, err
	}

	return &Playlist{
		ID:          playlist.ID,
		Name:        playlist.Name,
		ExternalURL: playlist.ExternalURLs.Spotify,
	}, nil
}

func AddTracksToPlaylist(accessToken, playlistID string, trackIDs []string) error {
	uris := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		uris[i] = "spotify:track:" + id
	}

	chunks := ChunkTrackIDs(uris, 100)

	for _, chunk := range chunks {
		url := fmt.Sprintf("%s/playlists/%s/tracks", APIURL, playlistID)

		body := addTracksRequest{URIs: chunk}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to add tracks: %d", resp.StatusCode)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func ClearPlaylist(accessToken, playlistID string) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", APIURL, playlistID)

	req, err := http.NewRequest("GET", url+"?limit=100", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tracksResp struct {
		Items []struct {
			Track struct {
				URI string `json:"uri"`
			} `json:"track"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tracksResp); err != nil {
		return err
	}

	if len(tracksResp.Items) == 0 {
		return nil
	}

	tracks := make([]map[string]string, len(tracksResp.Items))
	for i, item := range tracksResp.Items {
		tracks[i] = map[string]string{"uri": item.Track.URI}
	}

	deleteBody, _ := json.Marshal(map[string]interface{}{"tracks": tracks})

	req, err = http.NewRequest("DELETE", url, bytes.NewReader(deleteBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func GetUserPlaylists(accessToken string) ([]PlaylistItem, error) {
	var allPlaylists []PlaylistItem
	offset := 0
	limit := 50

	for {
		url := fmt.Sprintf("%s/me/playlists?limit=%d&offset=%d", APIURL, limit, offset)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get playlists: %d", resp.StatusCode)
		}

		var result PlaylistsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		allPlaylists = append(allPlaylists, result.Items...)

		if len(result.Items) < limit {
			break
		}
		offset += limit
	}

	return allPlaylists, nil
}

func FindExistingPlaylist(accessToken, playlistName string) (*Playlist, error) {
	url := fmt.Sprintf("%s/me/playlists?limit=50", APIURL)

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

	var playlistsResp struct {
		Items []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlistsResp); err != nil {
		return nil, err
	}

	for _, p := range playlistsResp.Items {
		if p.Name == playlistName {
			return &Playlist{
				ID:          p.ID,
				Name:        p.Name,
				ExternalURL: p.ExternalURLs.Spotify,
			}, nil
		}
	}

	return nil, nil
}
