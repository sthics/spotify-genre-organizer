const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://127.0.0.1:8080';

export async function fetchUser() {
  const res = await fetch(`${API_URL}/api/auth/me`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Not authenticated');
  }

  return res.json();
}

export async function startOrganize(playlistCount: number, replaceExisting: boolean) {
  const res = await fetch(`${API_URL}/api/organize`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      playlist_count: playlistCount,
      replace_existing: replaceExisting,
    }),
  });

  if (!res.ok) {
    throw new Error('Failed to start organize');
  }

  return res.json();
}

export async function getOrganizeStatus(jobId: string) {
  const res = await fetch(`${API_URL}/api/organize/${jobId}`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Failed to get status');
  }

  return res.json();
}

export async function logout() {
  await fetch(`${API_URL}/api/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  });
}

export async function getLibraryCount(): Promise<{ count: number; cached_at: string }> {
  const response = await fetch(`${API_URL}/api/library/count`, {
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Failed to fetch library count');
  }
  return response.json();
}

export interface UserSettings {
  user_id: string;
  name_template: string;
  description_template: string;
  is_premium: boolean;
}

export async function getSettings(): Promise<UserSettings> {
  const response = await fetch(`${API_URL}/api/settings`, {
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Failed to fetch settings');
  }
  return response.json();
}

export async function updateSettings(
  nameTemplate: string,
  descriptionTemplate: string
): Promise<UserSettings> {
  const response = await fetch(`${API_URL}/api/settings`, {
    method: 'PUT',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name_template: nameTemplate,
      description_template: descriptionTemplate,
    }),
  });
  if (!response.ok) {
    throw new Error('Failed to update settings');
  }
  return response.json();
}

export interface ManagedPlaylist {
  spotify_id: string;
  name: string;
  genre: string;
  song_count: number;
  spotify_url: string;
  image_url?: string;
  custom_name?: string;
  custom_description?: string;
  last_synced?: string;
}

export async function getPlaylists(): Promise<{
  playlists: ManagedPlaylist[];
  total_songs: number;
}> {
  const response = await fetch(`${API_URL}/api/playlists`, {
    credentials: 'include',
  });
  if (!response.ok) throw new Error('Failed to fetch playlists');
  return response.json();
}

export async function refreshPlaylist(id: string): Promise<{ song_count: number }> {
  const response = await fetch(`${API_URL}/api/playlists/${id}/refresh`, {
    method: 'POST',
    credentials: 'include',
  });
  if (!response.ok) throw new Error('Failed to refresh playlist');
  return response.json();
}

export async function updatePlaylist(
  id: string,
  customName?: string,
  customDescription?: string
): Promise<void> {
  const response = await fetch(`${API_URL}/api/playlists/${id}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      custom_name: customName,
      custom_description: customDescription,
    }),
  });
  if (!response.ok) throw new Error('Failed to update playlist');
}

export async function deletePlaylist(id: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/playlists/${id}`, {
    method: 'DELETE',
    credentials: 'include',
  });
  if (!response.ok) throw new Error('Failed to delete playlist');
}

export interface PlaylistSyncStatus {
  spotify_id: string;
  genre: string;
  new_count: number;
}

export interface SyncStatus {
  new_songs_count: number;
  oldest_sync_at: string | null;
  playlists: PlaylistSyncStatus[];
}

export interface SyncAllResult {
  playlists_updated: number;
  total_songs: number;
  failed_playlists?: string[];
}

// Custom error class for API errors with status codes
export class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message);
    this.name = 'ApiError';
  }
}

function handleApiResponse(response: Response): void {
  if (response.status === 401) {
    // Token expired - redirect to login
    window.location.href = '/';
    throw new ApiError('Session expired', 401);
  }
  if (!response.ok) {
    throw new ApiError(`Request failed: ${response.statusText}`, response.status);
  }
}

export async function getSyncStatus(): Promise<SyncStatus> {
  const response = await fetch(`${API_URL}/api/library/sync-status`, {
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

export async function syncAllPlaylists(): Promise<SyncAllResult> {
  const response = await fetch(`${API_URL}/api/playlists/sync-all`, {
    method: 'POST',
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

// Genre Analysis Types
export interface SubGenreCount {
  name: string;
  count: number;
}

export interface ParentGenreCount {
  name: string;
  count: number;
  sub_genres: SubGenreCount[];
}

export interface LibraryGenresResponse {
  parent_genres: ParentGenreCount[];
  analyzed_at: string;
  total_songs: number;
}

export async function getLibraryGenres(refresh = false): Promise<LibraryGenresResponse> {
  const url = refresh
    ? `${API_URL}/api/library/genres?refresh=true`
    : `${API_URL}/api/library/genres`;

  const response = await fetch(url, {
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

// Custom Playlist Types
export interface CustomPlaylistRequest {
  sub_genres: string[];
  mode: 'combined' | 'separate';
  name?: string;
}

export async function startCustomPlaylist(request: CustomPlaylistRequest): Promise<{ job_id: string }> {
  const response = await fetch(`${API_URL}/api/custom-playlist`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  handleApiResponse(response);
  return response.json();
}

// Exclusion Types
export interface ExclusionRule {
  id: string;
  user_id: string;
  exclusion_type: 'artist' | 'song';
  spotify_id: string;
  name: string;
  scope?: 'all_appearances' | 'primary_only';
  created_at: string;
}

export interface AffectedPlaylist {
  playlist_id: string;
  name: string;
  songs_to_remove: string[];
  song_count: number;
}

export interface ImpactPreview {
  total_songs_affected: number;
  affected_playlists: AffectedPlaylist[];
}

export interface CreateExclusionResponse {
  rule_id: string;
  preview: ImpactPreview | null;
}

export interface ApplyExclusionResponse {
  playlists_updated: number;
  songs_removed: number;
}

// Exclusion API Functions
export async function getExclusions(): Promise<{ rules: ExclusionRule[] }> {
  const response = await fetch(`${API_URL}/api/exclusions`, {
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

export async function createExclusion(
  type: 'artist' | 'song',
  spotifyId: string,
  name: string,
  scope?: 'all_appearances' | 'primary_only'
): Promise<CreateExclusionResponse> {
  const response = await fetch(`${API_URL}/api/exclusions`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      type,
      spotify_id: spotifyId,
      name,
      scope,
    }),
  });
  handleApiResponse(response);
  return response.json();
}

export async function applyExclusion(ruleId: string): Promise<ApplyExclusionResponse> {
  const response = await fetch(`${API_URL}/api/exclusions/${ruleId}/apply`, {
    method: 'POST',
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

export async function deleteExclusion(ruleId: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/exclusions/${ruleId}`, {
    method: 'DELETE',
    credentials: 'include',
  });
  handleApiResponse(response);
}
