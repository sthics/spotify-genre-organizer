'use client';

import { useState } from 'react';
import { Button } from '@/components/Button';
import { createExclusion, applyExclusion, CreateExclusionResponse } from '@/lib/api';

interface BlockArtistModalProps {
  artist: { id: string; name: string; imageUrl?: string };
  songCount?: number;
  isOpen: boolean;
  onClose: () => void;
  onBlocked: () => void;
}

type Step = 'scope' | 'preview' | 'loading';

export function BlockArtistModal({
  artist,
  songCount,
  isOpen,
  onClose,
  onBlocked,
}: BlockArtistModalProps) {
  const [step, setStep] = useState<Step>('scope');
  const [scope, setScope] = useState<'all_appearances' | 'primary_only'>('all_appearances');
  const [preview, setPreview] = useState<CreateExclusionResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleBlock = async () => {
    setStep('loading');
    setError(null);

    try {
      const result = await createExclusion('artist', artist.id, artist.name, scope);
      setPreview(result);

      if (result.preview && result.preview.total_songs_affected > 0) {
        setStep('preview');
      } else {
        // No songs affected, just close
        onBlocked();
        onClose();
      }
    } catch (err) {
      setError('Failed to create block rule. Please try again.');
      setStep('scope');
    }
  };

  const handleApply = async () => {
    if (!preview) return;

    setStep('loading');
    try {
      await applyExclusion(preview.rule_id);
      onBlocked();
      onClose();
    } catch (err) {
      setError('Failed to remove songs. Please try again.');
      setStep('preview');
    }
  };

  const handleSkip = () => {
    onBlocked();
    onClose();
  };

  const resetAndClose = () => {
    setStep('scope');
    setScope('all_appearances');
    setPreview(null);
    setError(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
      <div className="bg-bg-card rounded-2xl max-w-md w-full p-6 shadow-2xl animate-scale-in">
        {/* Header */}
        <div className="flex justify-between items-start mb-6">
          <div className="flex items-center gap-4">
            {artist.imageUrl ? (
              <img
                src={artist.imageUrl}
                alt={artist.name}
                className="w-12 h-12 rounded-full object-cover"
              />
            ) : (
              <div className="w-12 h-12 rounded-full bg-bg-dark flex items-center justify-center">
                <svg className="w-6 h-6 text-text-muted" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 14.016q2.906 0 4.945 2.039T19.031 21H4.969q0-2.906 2.039-4.945T12 14.016zm0-1.032q-1.641 0-2.813-1.172T8.015 9t1.172-2.813T12 5.015t2.813 1.172T15.985 9t-1.172 2.813T12 12.984z"/>
                </svg>
              </div>
            )}
            <div>
              <h2 className="font-display text-xl text-text-cream">
                {step === 'preview' ? `Blocking "${artist.name}"` : `Block "${artist.name}"?`}
              </h2>
              {songCount && step === 'scope' && (
                <p className="text-sm text-text-muted">{songCount} songs in your library</p>
              )}
            </div>
          </div>
          <button onClick={resetAndClose} className="text-text-muted hover:text-text-cream">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}

        {step === 'scope' && (
          <>
            <div className="space-y-3 mb-6">
              <button
                onClick={() => setScope('all_appearances')}
                className={`w-full p-4 rounded-xl border-2 text-left transition-all ${
                  scope === 'all_appearances'
                    ? 'border-accent-orange bg-accent-orange/10'
                    : 'border-bg-dark hover:border-text-muted'
                }`}
              >
                <div className="font-medium text-text-cream">All songs with this artist</div>
                <div className="text-sm text-text-muted">Including features & collabs</div>
              </button>
              <button
                onClick={() => setScope('primary_only')}
                className={`w-full p-4 rounded-xl border-2 text-left transition-all ${
                  scope === 'primary_only'
                    ? 'border-accent-orange bg-accent-orange/10'
                    : 'border-bg-dark hover:border-text-muted'
                }`}
              >
                <div className="font-medium text-text-cream">Only as main artist</div>
                <div className="text-sm text-text-muted">Keep their features</div>
              </button>
            </div>

            <div className="flex gap-3 justify-end">
              <Button variant="ghost" onClick={resetAndClose}>Cancel</Button>
              <Button
                onClick={handleBlock}
                className="bg-red-600 hover:bg-red-500"
              >
                Block Artist
              </Button>
            </div>
          </>
        )}

        {step === 'preview' && preview?.preview && (
          <>
            <p className="text-text-cream mb-4">
              This will remove {preview.preview.total_songs_affected} songs from{' '}
              {preview.preview.affected_playlists.length} playlists:
            </p>
            <div className="max-h-48 overflow-y-auto mb-4 space-y-2">
              {preview.preview.affected_playlists.map((p) => (
                <div key={p.playlist_id} className="flex justify-between items-center p-3 bg-bg-dark rounded-lg">
                  <span className="text-text-cream">{p.name}</span>
                  <span className="text-text-muted text-sm">{p.song_count} songs</span>
                </div>
              ))}
            </div>
            <p className="text-sm text-text-muted mb-6">
              Rule saved. Future playlists will also exclude this artist.
            </p>
            <div className="flex gap-3 justify-end">
              <Button variant="secondary" onClick={handleSkip}>Keep in playlists</Button>
              <Button onClick={handleApply}>Remove now</Button>
            </div>
          </>
        )}

        {step === 'loading' && (
          <div className="flex items-center justify-center py-8">
            <div className="w-8 h-8 border-2 border-accent-orange border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>
    </div>
  );
}
