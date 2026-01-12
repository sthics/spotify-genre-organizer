'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { PlaylistCard } from '@/components/PlaylistCard';
import { getOrganizeStatus } from '@/lib/api';

interface Playlist {
  name: string;
  genre: string;
  spotify_id: string;
  spotify_url: string;
  song_count: number;
}

function SuccessContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const jobId = searchParams.get('job');

  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [showConfetti, setShowConfetti] = useState(true);

  useEffect(() => {
    if (!jobId) {
      router.push('/dashboard');
      return;
    }

    getOrganizeStatus(jobId)
      .then((data) => {
        if (data.status === 'completed' && data.result) {
          setPlaylists(data.result.playlists);
        }
      })
      .catch(console.error);

    const timer = setTimeout(() => setShowConfetti(false), 3000);
    return () => clearTimeout(timer);
  }, [jobId, router]);

  const handleOpenAll = () => {
    if (playlists.length > 0) {
      window.open(playlists[0].spotify_url, '_blank');
    }
  };

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12 relative overflow-hidden">
      {/* Confetti particles */}
      {showConfetti && (
        <div className="fixed inset-0 pointer-events-none">
          {[...Array(20)].map((_, i) => (
            <div
              key={i}
              className="absolute w-2 h-2 rounded-full"
              style={{
                left: `${Math.random() * 100}%`,
                top: `-10px`,
                backgroundColor: ['#e85d04', '#f5f0e6', '#2d936c'][i % 3],
                animation: `fall ${3 + Math.random() * 2}s linear forwards`,
                animationDelay: `${Math.random() * 2}s`,
              }}
            />
          ))}
        </div>
      )}

      {/* Success Header */}
      <div className="text-center mb-12">
        <div className="text-4xl mb-4">&#10003;</div>
        <h1 className="font-display text-4xl text-text-cream mb-2">Done!</h1>
        <p className="text-text-muted text-xl">
          {playlists.length} playlists ready to play
        </p>
      </div>

      {/* Playlist Grid */}
      <div className="w-full max-w-4xl grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-12">
        {playlists.map((playlist, index) => (
          <PlaylistCard
            key={playlist.spotify_id}
            genre={playlist.genre}
            songCount={playlist.song_count}
            spotifyUrl={playlist.spotify_url}
            index={index}
          />
        ))}
      </div>

      {/* Actions */}
      <div className="flex flex-col sm:flex-row gap-4 items-center">
        <Button size="lg" onClick={handleOpenAll}>
          Open All in Spotify
          <span className="ml-2">&#9679;</span>
        </Button>
        <Button
          variant="ghost"
          onClick={() => router.push('/dashboard')}
        >
          Organize Again
        </Button>
      </div>

      {/* Add falling animation */}
      <style jsx>{`
        @keyframes fall {
          to {
            transform: translateY(100vh) rotate(720deg);
            opacity: 0;
          }
        }
      `}</style>
    </main>
  );
}

export default function Success() {
  return (
    <Suspense fallback={
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    }>
      <SuccessContent />
    </Suspense>
  );
}
