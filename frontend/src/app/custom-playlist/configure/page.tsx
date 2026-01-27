'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { useUser } from '@/hooks/useUser';
import { startCustomPlaylist } from '@/lib/api';

export default function ConfigurePage() {
    const { loading: userLoading } = useUser();
    const router = useRouter();
    const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
    const [mode, setMode] = useState<'combined' | 'separate'>('combined');
    const [playlistName, setPlaylistName] = useState('');
    const [isCreating, setIsCreating] = useState(false);

    useEffect(() => {
        // Get selected genres from sessionStorage
        const stored = sessionStorage.getItem('customPlaylistGenres');
        if (stored) {
            const genres = JSON.parse(stored);
            setSelectedGenres(genres);
            // If only one genre, skip to creation with that genre name
            if (genres.length === 1) {
                setMode('separate');
            }
        } else {
            // No genres selected, go back
            router.push('/custom-playlist');
        }
    }, [router]);

    const handleCreate = async () => {
        setIsCreating(true);
        try {
            const { job_id } = await startCustomPlaylist({
                sub_genres: selectedGenres,
                mode,
                name: mode === 'combined' ? (playlistName || 'Custom Mix by Organizer') : undefined,
            });
            router.push(`/processing?job=${job_id}`);
        } catch (error) {
            console.error('Failed to create playlist:', error);
            setIsCreating(false);
        }
    };

    if (userLoading || selectedGenres.length === 0) {
        return (
            <main className="min-h-screen flex items-center justify-center">
                <VinylIcon spinning size={64} />
            </main>
        );
    }

    // If only one genre, create immediately as separate (single playlist)
    if (selectedGenres.length === 1) {
        return (
            <main className="min-h-screen flex flex-col items-center justify-center px-4">
                <div className="max-w-md w-full bg-bg-card rounded-2xl p-8 text-center">
                    <h1 className="font-display text-2xl text-text-cream mb-4">
                        Creating {selectedGenres[0]} playlist
                    </h1>
                    <Button onClick={handleCreate} disabled={isCreating} className="w-full">
                        {isCreating ? (
                            <>
                                <VinylIcon spinning size={20} />
                                Creating...
                            </>
                        ) : (
                            'Create Playlist'
                        )}
                    </Button>
                </div>
            </main>
        );
    }

    return (
        <main className="min-h-screen flex flex-col items-center justify-center px-4 py-12">
            <div className="max-w-xl w-full">
                {/* Header */}
                <button
                    onClick={() => router.push('/custom-playlist')}
                    className="flex items-center gap-2 text-text-muted hover:text-text-cream transition-colors mb-6"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                    Back
                </button>

                <h1 className="font-display text-3xl text-text-cream mb-2">Almost there</h1>
                <div className="flex flex-wrap gap-2 mb-8">
                    {selectedGenres.map((genre) => (
                        <span key={genre} className="px-2 py-1 bg-bg-card text-text-cream text-sm rounded-full">
                            {genre}
                        </span>
                    ))}
                </div>

                {/* Options */}
                <div className="space-y-4">
                    {/* Combined Option */}
                    <button
                        onClick={() => setMode('combined')}
                        className={`w-full p-6 rounded-xl border-2 text-left transition-all ${mode === 'combined'
                                ? 'border-accent-orange bg-accent-orange/10'
                                : 'border-bg-card bg-bg-card hover:border-text-muted'
                            }`}
                    >
                        <div className="flex items-start gap-4">
                            <div className="p-3 bg-bg-dark rounded-lg">
                                <VinylIcon size={32} />
                            </div>
                            <div className="flex-1">
                                <h3 className="font-display text-xl text-text-cream mb-1">One Combined Playlist</h3>
                                <p className="text-text-muted text-sm mb-3">
                                    All songs from selected genres in one playlist
                                </p>
                                {mode === 'combined' && (
                                    <input
                                        type="text"
                                        value={playlistName}
                                        onChange={(e) => setPlaylistName(e.target.value)}
                                        placeholder="My Heavy Stuff"
                                        className="w-full px-3 py-2 bg-bg-dark border border-bg-card rounded-lg text-text-cream placeholder-text-muted focus:outline-none focus:border-accent-orange"
                                        onClick={(e) => e.stopPropagation()}
                                    />
                                )}
                            </div>
                        </div>
                    </button>

                    {/* Separate Option */}
                    <button
                        onClick={() => setMode('separate')}
                        className={`w-full p-6 rounded-xl border-2 text-left transition-all ${mode === 'separate'
                                ? 'border-accent-orange bg-accent-orange/10'
                                : 'border-bg-card bg-bg-card hover:border-text-muted'
                            }`}
                    >
                        <div className="flex items-start gap-4">
                            <div className="p-3 bg-bg-dark rounded-lg flex -space-x-2">
                                <VinylIcon size={24} />
                                <VinylIcon size={24} />
                                <VinylIcon size={24} />
                            </div>
                            <div className="flex-1">
                                <h3 className="font-display text-xl text-text-cream mb-1">Separate Playlists</h3>
                                <p className="text-text-muted text-sm">
                                    {selectedGenres.length} playlists, one per genre
                                </p>
                                <p className="text-text-muted text-xs mt-2">
                                    {selectedGenres.slice(0, 3).join(' • ')}{selectedGenres.length > 3 ? ' • ...' : ''}
                                </p>
                            </div>
                        </div>
                    </button>
                </div>

                {/* Create Button */}
                <div className="mt-8">
                    <Button onClick={handleCreate} disabled={isCreating} size="lg" className="w-full">
                        {isCreating ? (
                            <>
                                <VinylIcon spinning size={24} />
                                Creating...
                            </>
                        ) : (
                            <>
                                Create Playlist{mode === 'separate' && selectedGenres.length > 1 ? 's' : ''}
                                <span className="ml-2">●</span>
                            </>
                        )}
                    </Button>
                </div>
            </div>
        </main>
    );
}
