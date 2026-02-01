'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { GenrePicker } from '@/components/GenrePicker';
import { Skeleton } from '@/components/Skeleton';
import { useUser } from '@/hooks/useUser';
import { getLibraryGenres, LibraryGenresResponse } from '@/lib/api';

const LOADING_MESSAGES = [
    "Dusting off the vinyls...",
    "Categorizing the vibes...",
    "Finding your hidden gems...",
    "Tuning the frequencies...",
    "Digging deep into the crates...",
    "Polishing the grooves..."
];

export default function CustomPlaylistPage() {
    const { user, loading: userLoading } = useUser();
    const router = useRouter();
    const [genreData, setGenreData] = useState<LibraryGenresResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [selectedGenres, setSelectedGenres] = useState<Set<string>>(new Set());
    const [loadingMessageIndex, setLoadingMessageIndex] = useState(0);

    useEffect(() => {
        const fetchGenres = async () => {
            try {
                const data = await getLibraryGenres();
                setGenreData(data);
            } catch (error) {
                console.error('Failed to fetch genres:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchGenres();
    }, []);

    // Rotate loading messages
    useEffect(() => {
        if (!loading) return;
        const interval = setInterval(() => {
            setLoadingMessageIndex((prev) => (prev + 1) % LOADING_MESSAGES.length);
        }, 2000);
        return () => clearInterval(interval);
    }, [loading]);

    const handleRefresh = async () => {
        setRefreshing(true);
        try {
            const data = await getLibraryGenres(true);
            setGenreData(data);
        } catch (error) {
            console.error('Failed to refresh genres:', error);
        } finally {
            setRefreshing(false);
        }
    };

    const getTotalSelectedSongs = () => {
        if (!genreData) return 0;
        let total = 0;
        for (const parent of genreData.parent_genres) {
            for (const sub of parent.sub_genres) {
                if (selectedGenres.has(sub.name)) {
                    total += sub.count;
                }
            }
        }
        return total;
    };

    const formatAnalyzedTime = (dateStr: string) => {
        const date = new Date(dateStr);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'just now';
        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;
        return `${diffDays}d ago`;
    };

    const handleNext = () => {
        // Store selection in sessionStorage and navigate to configure page
        sessionStorage.setItem('customPlaylistGenres', JSON.stringify(Array.from(selectedGenres)));
        router.push('/custom-playlist/configure');
    };

    // Show initial spinner only for auth check, or skeleton for data load
    if (userLoading) {
        return (
            <main className="min-h-screen flex items-center justify-center">
                <VinylIcon spinning size={64} />
            </main>
        );
    }

    if (loading) {
        return (
            <main className="min-h-screen flex flex-col px-4 py-8">
                {/* Skeleton Header */}
                <div className="max-w-2xl mx-auto w-full mb-8">
                    <div className="flex items-center gap-2 mb-4">
                        <Skeleton className="w-5 h-5 rounded-full" />
                        <Skeleton className="w-12 h-5" />
                    </div>

                    <h1 className="font-display text-3xl text-text-cream mb-2 animate-pulse">
                        Digging through {user?.display_name?.split(' ')[0] || 'your'}&apos;s crates...
                    </h1>
                    <p className="text-text-muted italic animate-fade-in key={loadingMessageIndex}">
                        {LOADING_MESSAGES[loadingMessageIndex]}
                    </p>
                </div>

                {/* Skeleton Genre Grid */}
                <div className="max-w-2xl mx-auto w-full flex-1 mb-24 space-y-2">
                    {[...Array(8)].map((_, i) => (
                        <Skeleton key={i} className="w-full h-20 rounded-xl" />
                    ))}
                </div>

                {/* Skeleton Bottom Bar */}
                <div className="fixed bottom-0 left-0 right-0 bg-bg-dark border-t border-bg-card p-4">
                    <div className="max-w-2xl mx-auto flex items-center justify-between">
                        <Skeleton className="w-48 h-8" />
                        <Skeleton className="w-24 h-10 rounded-lg" />
                    </div>
                </div>
            </main>
        );
    }

    return (
        <main className="min-h-screen flex flex-col px-4 py-8">
            {/* Header */}
            <div className="max-w-2xl mx-auto w-full mb-8">
                <button
                    onClick={() => router.push('/dashboard')}
                    className="flex items-center gap-2 text-text-muted hover:text-text-cream transition-colors mb-4"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                    Back
                </button>

                <h1 className="font-display text-3xl text-text-cream mb-2">Build a Crate</h1>

                {genreData && (
                    <div className="flex items-center gap-2 text-text-muted">
                        <span>{genreData.total_songs.toLocaleString()} songs</span>
                        <span>•</span>
                        <span>Analyzed {formatAnalyzedTime(genreData.analyzed_at)}</span>
                        <button
                            onClick={handleRefresh}
                            disabled={refreshing}
                            className="ml-1 p-1 hover:text-text-cream transition-colors disabled:opacity-50"
                            title="Refresh analysis"
                        >
                            <svg
                                className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`}
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                            </svg>
                        </button>
                    </div>
                )}
            </div>

            {/* Genre Picker */}
            <div className="max-w-2xl mx-auto w-full flex-1 mb-24">
                {genreData && (
                    <GenrePicker
                        parentGenres={genreData.parent_genres}
                        selectedGenres={selectedGenres}
                        onSelectionChange={setSelectedGenres}
                    />
                )}
            </div>

            {/* Sticky Bottom Bar */}
            <div className="fixed bottom-0 left-0 right-0 bg-bg-dark border-t border-bg-card p-4">
                <div className="max-w-2xl mx-auto flex items-center justify-between">
                    <div className="flex items-center gap-2 flex-wrap flex-1 mr-4">
                        {selectedGenres.size === 0 ? (
                            <span className="text-text-muted">Select genres to continue</span>
                        ) : (
                            <>
                                {Array.from(selectedGenres).slice(0, 3).map((genre) => (
                                    <span
                                        key={genre}
                                        className="px-2 py-1 bg-accent-orange text-white text-sm rounded-full flex items-center gap-1 animate-bounce-in"
                                    >
                                        {genre}
                                        <button
                                            onClick={() => {
                                                const newSelected = new Set(selectedGenres);
                                                newSelected.delete(genre);
                                                setSelectedGenres(newSelected);
                                            }}
                                            className="hover:bg-white/20 rounded-full p-0.5"
                                        >
                                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                        </button>
                                    </span>
                                ))}
                                {selectedGenres.size > 3 && (
                                    <span className="text-text-muted text-sm">+{selectedGenres.size - 3} more</span>
                                )}
                            </>
                        )}
                    </div>
                    <div className="flex items-center gap-4">
                        <span className="text-text-cream font-medium">
                            {getTotalSelectedSongs().toLocaleString()} songs
                        </span>
                        <Button
                            onClick={handleNext}
                            disabled={selectedGenres.size === 0}
                        >
                            Next
                            <svg className="w-4 h-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                            </svg>
                        </Button>
                    </div>
                </div>
            </div>
        </main>
    );
}
