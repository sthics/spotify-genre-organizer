package organizer

import (
	"log"
	"sort"
	"time"

	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type OrganizeResult struct {
	Playlists []PlaylistResult `json:"playlists"`
}

type PlaylistResult struct {
	Name       string `json:"name"`
	Genre      string `json:"genre"`
	SpotifyID  string `json:"spotify_id"`
	SpotifyURL string `json:"spotify_url"`
	SongCount  int    `json:"song_count"`
}

type ProgressCallback func(stage string, processed, total int)

func OrganizeSongs(
	accessToken string,
	userID string,
	songs []spotify.Song,
	playlistCount int,
	replaceExisting bool,
	progress ProgressCallback,
) (*OrganizeResult, error) {
	// Fetch user settings
	settings, err := database.GetUserSettings(userID)
	if err != nil {
		log.Printf("Failed to fetch settings for user %s: %v", userID, err)
		settings = models.DefaultSettings(userID)
	}

	// Group songs by parent genre
	genreGroups := make(map[string][]spotify.Song)
	for _, song := range songs {
		if len(song.Genres) == 0 {
			genreGroups["Other"] = append(genreGroups["Other"], song)
			continue
		}

		// Use weighted scoring across all genres
		parentGenre := genres.ScoreGenres(song.Genres)
		genreGroups[parentGenre] = append(genreGroups[parentGenre], song)
	}

	// Sort genres by song count (descending)
	type genreCount struct {
		genre string
		count int
	}
	var sortedGenres []genreCount
	for genre, songs := range genreGroups {
		sortedGenres = append(sortedGenres, genreCount{genre, len(songs)})
	}
	sort.Slice(sortedGenres, func(i, j int) bool {
		return sortedGenres[i].count > sortedGenres[j].count
	})

	// Limit to requested playlist count
	if len(sortedGenres) > playlistCount {
		// Merge smaller genres into "Other"
		for i := playlistCount; i < len(sortedGenres); i++ {
			genreGroups["Other"] = append(genreGroups["Other"], genreGroups[sortedGenres[i].genre]...)
			delete(genreGroups, sortedGenres[i].genre)
		}
		sortedGenres = sortedGenres[:playlistCount]
	}

	// Create playlists
	var results []PlaylistResult
	total := len(sortedGenres)

	for i, gc := range sortedGenres {
		if progress != nil {
			progress("creating", i+1, total)
		}

		// Create or Update Playlist
		playlistName := settings.BuildPlaylistName(gc.genre)
		playlistDescription := settings.BuildDescription(gc.genre)
		songs := genreGroups[gc.genre]

		var playlist *spotify.Playlist
		var err error

		if replaceExisting {
			// Check for existing playlist
			playlist, err = spotify.FindExistingPlaylist(accessToken, playlistName)
			if err != nil {
				return nil, err
			}

			if playlist != nil {
				// Clear existing tracks
				if err := spotify.ClearPlaylist(accessToken, playlist.ID); err != nil {
					return nil, err
				}
			}
		}

		if playlist == nil {
			// Create new playlist
			playlist, err = spotify.CreatePlaylist(
				accessToken,
				userID,
				playlistName,
				playlistDescription,
			)
			if err != nil {
				return nil, err
			}
		}

		// Add tracks
		trackIDs := make([]string, len(songs))
		for i, s := range songs {
			trackIDs[i] = s.ID
		}

		if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, trackIDs); err != nil {
			return nil, err
		}

		// Save playlist override with last_synced_at for sync tracking
		now := time.Now()
		override := &models.PlaylistOverride{
			UserID:            userID,
			PlaylistSpotifyID: playlist.ID,
			Genre:             gc.genre,
			LastSyncedAt:      &now,
		}
		if err := database.SavePlaylistOverride(override); err != nil {
			log.Printf("Failed to save playlist override for %s: %v", playlist.ID, err)
			// Don't fail the whole operation for this
		}

		results = append(results, PlaylistResult{
			Name:       playlistName,
			Genre:      gc.genre,
			SpotifyID:  playlist.ID,
			SpotifyURL: playlist.ExternalURL,
			SongCount:  len(songs),
		})
	}

	return &OrganizeResult{Playlists: results}, nil
}
