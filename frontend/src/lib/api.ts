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
