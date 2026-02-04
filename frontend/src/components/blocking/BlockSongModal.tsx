'use client';

import { useState } from 'react';
import { Button } from '@/components/Button';
import { createExclusion, applyExclusion, CreateExclusionResponse } from '@/lib/api';

interface BlockSongModalProps {
  song: { id: string; name: string; artistName: string };
  isOpen: boolean;
  onClose: () => void;
  onBlocked: () => void;
}

type Step = 'confirm' | 'preview' | 'loading';

export function BlockSongModal({
  song,
  isOpen,
  onClose,
  onBlocked,
}: BlockSongModalProps) {
  const [step, setStep] = useState<Step>('confirm');
  const [preview, setPreview] = useState<CreateExclusionResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleBlock = async () => {
    setStep('loading');
    setError(null);

    try {
      const result = await createExclusion('song', song.id, song.name);
      setPreview(result);

      if (result.preview && result.preview.total_songs_affected > 0) {
        setStep('preview');
      } else {
        onBlocked();
        onClose();
      }
    } catch (err) {
      setError('Failed to block song. Please try again.');
      setStep('confirm');
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
      setError('Failed to remove song. Please try again.');
      setStep('preview');
    }
  };

  const handleSkip = () => {
    onBlocked();
    onClose();
  };

  const resetAndClose = () => {
    setStep('confirm');
    setPreview(null);
    setError(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
      <div className="bg-bg-card rounded-2xl max-w-md w-full p-6 shadow-2xl animate-scale-in">
        <div className="flex justify-between items-start mb-6">
          <div>
            <h2 className="font-display text-xl text-text-cream">Block "{song.name}"?</h2>
            <p className="text-sm text-text-muted">{song.artistName}</p>
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

        {step === 'confirm' && (
          <>
            <p className="text-text-muted mb-6">
              This song will be removed from all playlists and never added again.
            </p>
            <div className="flex gap-3 justify-end">
              <Button variant="ghost" onClick={resetAndClose}>Cancel</Button>
              <Button onClick={handleBlock} className="bg-red-600 hover:bg-red-500">
                Block Song
              </Button>
            </div>
          </>
        )}

        {step === 'preview' && preview?.preview && (
          <>
            <p className="text-text-cream mb-4">
              This will remove the song from {preview.preview.affected_playlists.length} playlists.
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
