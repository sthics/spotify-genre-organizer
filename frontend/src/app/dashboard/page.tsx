'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { Slider } from '@/components/Slider';
import { useUser } from '@/hooks/useUser';
import { startOrganize, logout, getLibraryCount } from '@/lib/api';

const LIBRARY_COUNT_CACHE_KEY = 'spotify_library_count';

export default function Dashboard() {
  const { user, loading } = useUser();
  const router = useRouter();
  const [playlistCount, setPlaylistCount] = useState(12);
  const [replaceExisting, setReplaceExisting] = useState(true);
  const [isOrganizing, setIsOrganizing] = useState(false);
  const [likedSongsCount, setLikedSongsCount] = useState<number | null>(null);
  const [countLoading, setCountLoading] = useState(true);

  useEffect(() => {
    // Check localStorage first for instant display
    const cached = localStorage.getItem(LIBRARY_COUNT_CACHE_KEY);
    if (cached) {
      try {
        const { count } = JSON.parse(cached);
        setLikedSongsCount(count);
      } catch {
        // Invalid cache, ignore
      }
    }

    // Fetch fresh count from API
    const fetchCount = async () => {
      try {
        const data = await getLibraryCount();
        setLikedSongsCount(data.count);
        localStorage.setItem(LIBRARY_COUNT_CACHE_KEY, JSON.stringify(data));
      } catch (error) {
        console.error('Failed to fetch library count:', error);
      } finally {
        setCountLoading(false);
      }
    };

    fetchCount();
  }, []);

  const songsPerPlaylist = likedSongsCount ? Math.round(likedSongsCount / playlistCount) : 0;

  const handleOrganize = async () => {
    setIsOrganizing(true);
    try {
      const { job_id } = await startOrganize(playlistCount, replaceExisting);
      router.push(`/processing?job=${job_id}`);
    } catch (error) {
      console.error('Failed to start organize:', error);
      setIsOrganizing(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    router.push('/');
  };

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4 py-12">
      {/* Header */}
      <div className="w-full max-w-xl flex items-center justify-between mb-12">
        <div>
          <h1 className="font-display text-2xl text-text-cream">
            Hey, {user?.display_name?.split(' ')[0] || 'there'}.
          </h1>
          <p className="text-text-muted">
            You&apos;ve got{' '}
            {countLoading && likedSongsCount === null ? (
              <span className="inline-block w-16 h-5 bg-bg-card rounded animate-pulse align-middle" />
            ) : (
              <span className="text-text-cream font-medium">{(likedSongsCount ?? 0).toLocaleString()}</span>
            )}{' '}
            liked songs.
          </p>
        </div>
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.push('/settings')}
            className="text-text-muted hover:text-text-cream transition-colors"
            title="Settings"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
          <button
            onClick={handleLogout}
            className="text-text-muted hover:text-text-cream transition-colors"
            title="Logout"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </div>

      {/* Main Card */}
      <div className="w-full max-w-xl bg-bg-card rounded-2xl p-8 shadow-xl">
        {/* Playlist Count Slider */}
        <div className="mb-8">
          <h2 className="font-display text-xl text-text-cream mb-4">
            How many playlists?
          </h2>
          <Slider
            value={playlistCount}
            onChange={setPlaylistCount}
            min={1}
            max={50}
          />
          <div className="flex items-center gap-2 mt-4 text-text-muted">
            <VinylIcon size={20} />
            <span>~{songsPerPlaylist} songs per playlist</span>
          </div>
        </div>

        {/* Replace Toggle */}
        <div className="mb-8">
          <div className="space-y-3">
            <label className="flex items-start gap-3 cursor-pointer group">
              <input
                type="radio"
                name="replace"
                checked={replaceExisting}
                onChange={() => setReplaceExisting(true)}
                className="mt-1 accent-accent-orange"
              />
              <div>
                <span className="text-text-cream group-hover:text-accent-orange transition-colors">
                  Update existing playlists
                </span>
                <p className="text-sm text-text-muted">
                  Replaces songs in &quot;Rock by Organizer&quot;, etc.
                </p>
              </div>
            </label>
            <label className="flex items-start gap-3 cursor-pointer group">
              <input
                type="radio"
                name="replace"
                checked={!replaceExisting}
                onChange={() => setReplaceExisting(false)}
                className="mt-1 accent-accent-orange"
              />
              <div>
                <span className="text-text-cream group-hover:text-accent-orange transition-colors">
                  Create fresh playlists
                </span>
                <p className="text-sm text-text-muted">
                  Keeps your old ones, makes new
                </p>
              </div>
            </label>
          </div>
        </div>

        {/* Manage Playlists Link */}
        <div className="mb-4">
          <Button
            size="lg"
            variant="secondary"
            className="w-full flex items-center justify-center gap-2"
            onClick={() => router.push('/playlists')}
          >
            <VinylIcon size={20} />
            Manage My Crates
          </Button>
        </div>

        {/* Organize Button */}
        <Button
          size="lg"
          className="w-full flex items-center justify-center gap-2"
          onClick={handleOrganize}
          disabled={isOrganizing}
        >
          {isOrganizing ? (
            <>
              <VinylIcon spinning size={24} />
              Starting...
            </>
          ) : (
            <>
              Organize My Library
              <span className="text-xl">&#9679;</span>
            </>
          )}
        </Button>
      </div>
    </main>
  );
}
