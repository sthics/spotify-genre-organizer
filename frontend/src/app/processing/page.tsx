'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { ProgressBar } from '@/components/ProgressBar';
import { GenreTag } from '@/components/GenreTag';
import { getOrganizeStatus } from '@/lib/api';

interface JobStatus {
  id: string;
  status: string;
  stage: string;
  songs_processed: number;
  total_songs: number;
  genres_discovered: string[];
  result?: {
    playlists: Array<{
      name: string;
      genre: string;
      spotify_id: string;
      spotify_url: string;
      song_count: number;
    }>;
  };
  error?: string;
}

function ProcessingContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const jobId = searchParams.get('job');

  const [status, setStatus] = useState<JobStatus | null>(null);
  const [tonearmAngle, setTonearmAngle] = useState(0);

  useEffect(() => {
    if (!jobId) {
      router.push('/dashboard');
      return;
    }

    const pollStatus = async () => {
      try {
        const data = await getOrganizeStatus(jobId);
        setStatus(data);

        if (data.total_songs > 0) {
          const progress = data.songs_processed / data.total_songs;
          setTonearmAngle(progress * 30);
        }

        if (data.status === 'completed') {
          router.push(`/success?job=${jobId}`);
        } else if (data.status === 'failed') {
          console.error('Job failed:', data.error);
        } else {
          setTimeout(pollStatus, 1000);
        }
      } catch (error) {
        console.error('Failed to get status:', error);
        setTimeout(pollStatus, 2000);
      }
    };

    pollStatus();
  }, [jobId, router]);

  const getStageText = (stage: string) => {
    switch (stage) {
      case 'fetching':
        return 'Analyzing your library...';
      case 'analyzing':
        return 'Detecting genres...';
      case 'creating':
        return 'Creating playlists...';
      default:
        return 'Processing...';
    }
  };

  const progress = status?.total_songs
    ? (status.songs_processed / status.total_songs) * 100
    : 0;

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4">
      {/* Vinyl with Tonearm */}
      <div className="relative mb-8">
        <VinylIcon spinning size={120} />
        {/* Tonearm */}
        <div
          className="absolute top-0 right-0 w-16 h-1 bg-text-muted origin-right transition-transform duration-500"
          style={{
            transform: `rotate(${-45 + tonearmAngle}deg)`,
            transformOrigin: 'right center',
          }}
        >
          <div className="absolute right-0 top-1/2 -translate-y-1/2 w-3 h-3 bg-text-cream rounded-full" />
        </div>
      </div>

      {/* Status Text */}
      <h2 className="font-display text-2xl text-text-cream mb-4">
        {getStageText(status?.stage || '')}
      </h2>

      {/* Progress Bar */}
      <div className="w-full max-w-md mb-8">
        <ProgressBar
          progress={progress}
          label={status?.total_songs
            ? `${status.songs_processed.toLocaleString()} / ${status.total_songs.toLocaleString()} songs`
            : undefined
          }
        />
      </div>

      {/* Discovered Genres */}
      {status?.genres_discovered && status.genres_discovered.length > 0 && (
        <div className="w-full max-w-lg">
          <p className="text-text-muted text-center mb-4">Genres discovered:</p>
          <div className="flex flex-wrap gap-2 justify-center">
            {status.genres_discovered.slice(0, 12).map((genre, index) => (
              <GenreTag key={genre} genre={genre} index={index} />
            ))}
            {status.genres_discovered.length > 12 && (
              <span className="text-text-muted">
                +{status.genres_discovered.length - 12} more
              </span>
            )}
          </div>
        </div>
      )}
    </main>
  );
}

export default function Processing() {
  return (
    <Suspense fallback={
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    }>
      <ProcessingContent />
    </Suspense>
  );
}
