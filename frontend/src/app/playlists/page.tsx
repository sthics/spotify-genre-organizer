'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import {
    getPlaylists,
    refreshPlaylist,
    updatePlaylist,
    deletePlaylist,
    ManagedPlaylist,
} from '@/lib/api';

const GENRE_COLORS: Record<string, string> = {
    Rock: '#e85d04',
    Electronic: '#00d4ff',
    'Hip-Hop': '#ffd700',
    Jazz: '#9b59b6',
    Reggae: '#2d936c',
    Pop: '#ff69b4',
    'R&B': '#ff6b6b',
    Metal: '#8b0000',
    Folk: '#8b4513',
    Classical: '#daa520',
    Other: '#888888',
};

export default function Playlists() {
    const router = useRouter();
    const [playlists, setPlaylists] = useState<ManagedPlaylist[]>([]);
    const [totalSongs, setTotalSongs] = useState(0);
    const [loading, setLoading] = useState(true);
    const [expandedId, setExpandedId] = useState<string | null>(null);
    const [refreshingId, setRefreshingId] = useState<string | null>(null);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editName, setEditName] = useState('');
    const [editDesc, setEditDesc] = useState('');

    useEffect(() => {
        loadPlaylists();
    }, []);

    const loadPlaylists = async () => {
        try {
            const data = await getPlaylists();
            setPlaylists(data.playlists || []);
            setTotalSongs(data.total_songs);
        } catch (err) {
            console.error(err);
        }
        setLoading(false);
    };

    const handleRefresh = async (id: string) => {
        setRefreshingId(id);
        try {
            const result = await refreshPlaylist(id);
            setPlaylists((prev) =>
                prev.map((p) =>
                    p.spotify_id === id ? { ...p, song_count: result.song_count } : p
                )
            );
        } catch (err) {
            console.error(err);
        }
        setRefreshingId(null);
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Delete this playlist from Spotify?')) return;
        try {
            await deletePlaylist(id);
            setPlaylists((prev) => prev.filter((p) => p.spotify_id !== id));
            if (expandedId === id) setExpandedId(null);
        } catch (err) {
            console.error(err);
            alert('Failed to delete playlist');
        }
    };

    const startEditing = (playlist: ManagedPlaylist) => {
        setEditingId(playlist.spotify_id);
        setEditName(playlist.custom_name || playlist.name);
        setEditDesc(playlist.custom_description || '');
    };

    const saveEdit = async (id: string) => {
        try {
            await updatePlaylist(id, editName, editDesc);
            setPlaylists((prev) =>
                prev.map((p) =>
                    p.spotify_id === id
                        ? { ...p, custom_name: editName, custom_description: editDesc }
                        : p
                )
            );
            setEditingId(null);
        } catch (err) {
            console.error(err);
            alert('Failed to save changes');
        }
    };

    const toggleExpand = (id: string) => {
        if (expandedId === id) {
            setExpandedId(null);
            setEditingId(null);
        } else {
            setExpandedId(id);
            setEditingId(null);
        }
    };

    if (loading) {
        return (
            <main className="min-h-screen flex items-center justify-center">
                <VinylIcon spinning size={64} />
            </main>
        );
    }

    return (
        <main className="min-h-screen flex flex-col items-center px-4 py-12">
            <div className="w-full max-w-5xl">
                {/* Header */}
                <div className="flex items-center justify-between mb-12">
                    <div className="flex items-center gap-4">
                        <button
                            onClick={() => router.push('/dashboard')}
                            className="text-text-muted hover:text-text-cream transition-colors"
                        >
                            &larr; Back
                        </button>
                        <div>
                            <h1 className="font-display text-2xl text-text-cream flex items-center gap-2">
                                Your Crates
                            </h1>
                            <p className="text-text-muted">
                                {playlists.length} playlists &bull; {totalSongs} songs organized
                            </p>
                        </div>
                    </div>
                    <Button onClick={() => router.push('/dashboard')}>
                        + New Organize
                    </Button>
                </div>

                {playlists.length === 0 ? (
                    <div className="text-center py-20 text-text-muted">
                        <p>No organized playlists found.</p>
                        <Button
                            className="mt-4"
                            variant="secondary"
                            onClick={() => router.push('/dashboard')}
                        >
                            Start Organizing
                        </Button>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 px-4">
                        {playlists.map((playlist) => {
                            const isExpanded = expandedId === playlist.spotify_id;
                            const isRefreshing = refreshingId === playlist.spotify_id;
                            const isEditing = editingId === playlist.spotify_id;
                            const color = GENRE_COLORS[playlist.genre] || GENRE_COLORS['Other'];

                            return (
                                <div
                                    key={playlist.spotify_id}
                                    className={`relative group perspective-1000 transition-all duration-300 ${isExpanded ? 'z-10 md:col-span-2 lg:col-span-3' : 'z-0'
                                        }`}
                                >
                                    {/* Record Sleeve Card */}
                                    <div
                                        onClick={() => !isEditing && toggleExpand(playlist.spotify_id)}
                                        className={`
                      relative bg-bg-card rounded-lg shadow-xl cursor-pointer
                      transition-all duration-500 transform preserve-3d
                      ${isExpanded ? 'rotate-x-0 bg-bg-dark border border-gray-700' : 'hover:-translate-y-2 hover:rotate-x-2'}
                    `}
                                        style={{
                                            borderLeft: `4px solid ${color}`,
                                        }}
                                    >
                                        <div className="p-6">
                                            <div className="flex justify-between items-start mb-4">
                                                <div>
                                                    <span
                                                        className="inline-block px-2 py-0.5 rounded text-xs font-bold uppercase tracking-wider mb-2"
                                                        style={{ backgroundColor: `${color}33`, color: color }}
                                                    >
                                                        {playlist.genre}
                                                    </span>

                                                    {isEditing ? (
                                                        <div className="mt-2 space-y-3" onClick={(e) => e.stopPropagation()}>
                                                            <input
                                                                type="text"
                                                                value={editName}
                                                                onChange={(e) => setEditName(e.target.value)}
                                                                className="w-full bg-bg-dark border border-gray-600 rounded px-3 py-2 text-text-cream"
                                                                placeholder="Playlist Name"
                                                                autoFocus
                                                            />
                                                            <textarea
                                                                value={editDesc}
                                                                onChange={(e) => setEditDesc(e.target.value)}
                                                                className="w-full bg-bg-dark border border-gray-600 rounded px-3 py-2 text-text-cream text-sm"
                                                                placeholder="Description"
                                                                rows={2}
                                                            />
                                                            <div className="flex gap-2">
                                                                <Button size="sm" onClick={() => saveEdit(playlist.spotify_id)}>
                                                                    Save
                                                                </Button>
                                                                <Button
                                                                    size="sm"
                                                                    variant="secondary"
                                                                    onClick={() => setEditingId(null)}
                                                                >
                                                                    Cancel
                                                                </Button>
                                                            </div>
                                                        </div>
                                                    ) : (
                                                        <>
                                                            <h3 className="font-display text-xl text-text-cream mb-1">
                                                                {playlist.custom_name || playlist.name}
                                                            </h3>
                                                            {playlist.song_count > 0 && (
                                                                <p className="text-sm text-text-muted">
                                                                    {playlist.song_count} songs
                                                                </p>
                                                            )}
                                                        </>
                                                    )}
                                                </div>

                                                {isRefreshing ? (
                                                    <VinylIcon spinning size={24} />
                                                ) : (
                                                    <div className={`transition-transform duration-300 ${isExpanded ? 'rotate-180' : ''}`}>
                                                        <svg className="w-6 h-6 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                                        </svg>
                                                    </div>
                                                )}
                                            </div>

                                            {/* Expanded Actions */}
                                            <div
                                                className={`overflow-hidden transition-all duration-300 ${isExpanded ? 'max-h-48 opacity-100 mt-4' : 'max-h-0 opacity-0'
                                                    }`}
                                            >
                                                <div className="h-px bg-gray-700 mb-4" />

                                                <div className="flex flex-wrap gap-3" onClick={(e) => e.stopPropagation()}>
                                                    <Button
                                                        size="sm"
                                                        variant="secondary"
                                                        onClick={() => handleRefresh(playlist.spotify_id)}
                                                        disabled={isRefreshing}
                                                    >
                                                        {isRefreshing ? 'Syncing...' : '↻ Refresh'}
                                                    </Button>

                                                    <Button
                                                        size="sm"
                                                        variant="secondary"
                                                        onClick={() => startEditing(playlist)}
                                                    >
                                                        ✎ Edit Details
                                                    </Button>

                                                    <a
                                                        href={playlist.spotify_url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="inline-flex items-center justify-center px-4 py-2 rounded-full 
                                      bg-gray-700 text-text-cream text-sm font-medium hover:bg-gray-600 transition-colors"
                                                    >
                                                        Open in Spotify ↗
                                                    </a>

                                                    <div className="flex-grow" />

                                                    <button
                                                        onClick={() => handleDelete(playlist.spotify_id)}
                                                        className="text-red-400 hover:text-red-300 text-sm font-medium px-3"
                                                    >
                                                        🗑 Delete
                                                    </button>
                                                </div>

                                                {playlist.last_synced && (
                                                    <p className="text-xs text-text-muted mt-4 text-right">
                                                        Last synced: {new Date(playlist.last_synced).toLocaleDateString()}
                                                    </p>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </main>
    );
}
