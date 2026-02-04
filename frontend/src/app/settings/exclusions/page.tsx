'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/Button';
import { getExclusions, deleteExclusion, ExclusionRule } from '@/lib/api';
import { useToast } from '@/contexts/ToastContext';

export default function ExclusionsSettings() {
  const router = useRouter();
  const { showToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [rules, setRules] = useState<ExclusionRule[]>([]);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  useEffect(() => {
    loadRules();
  }, []);

  const loadRules = async () => {
    try {
      const data = await getExclusions();
      setRules(data.rules || []);
    } catch (err) {
      showToast('Failed to load exclusions', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (rule: ExclusionRule) => {
    if (!confirm(`Unblock "${rule.name}"? Their songs can appear in your playlists again.`)) {
      return;
    }

    setDeletingId(rule.id);
    try {
      await deleteExclusion(rule.id);
      setRules(rules.filter((r) => r.id !== rule.id));
      showToast(`${rule.name} unblocked`, 'success');
    } catch (err) {
      showToast('Failed to unblock', 'error');
    } finally {
      setDeletingId(null);
    }
  };

  const artistRules = rules.filter((r) => r.exclusion_type === 'artist');
  const songRules = rules.filter((r) => r.exclusion_type === 'song');

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <div className="text-text-cream animate-pulse">Loading exclusions...</div>
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12">
      <div className="w-full max-w-2xl">
        <button
          onClick={() => router.push('/settings')}
          className="text-text-muted hover:text-text-cream transition-colors mb-8"
        >
          &larr; Back to Settings
        </button>

        <h1 className="font-display text-3xl text-text-cream mb-8">Blocked Artists & Songs</h1>

        {/* Artists Section */}
        <section className="mb-8">
          <h2 className="font-display text-xl text-text-cream mb-4">Blocked Artists</h2>
          {artistRules.length === 0 ? (
            <div className="bg-bg-card rounded-xl p-6 text-center text-text-muted">
              No blocked artists yet. Block artists while browsing your library.
            </div>
          ) : (
            <div className="bg-bg-card rounded-xl overflow-hidden divide-y divide-bg-dark">
              {artistRules.map((rule) => (
                <div key={rule.id} className="flex items-center justify-between p-4">
                  <div>
                    <div className="font-medium text-text-cream">{rule.name}</div>
                    <div className="text-sm text-text-muted">
                      {rule.scope === 'all_appearances' ? 'All appearances' : 'Main artist only'}
                    </div>
                  </div>
                  <button
                    onClick={() => handleDelete(rule)}
                    disabled={deletingId === rule.id}
                    className="p-2 text-text-muted hover:text-red-400 transition-colors disabled:opacity-50"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Songs Section */}
        <section>
          <h2 className="font-display text-xl text-text-cream mb-4">Blocked Songs</h2>
          {songRules.length === 0 ? (
            <div className="bg-bg-card rounded-xl p-6 text-center text-text-muted">
              No blocked songs yet. Block songs while browsing your library or playlists.
            </div>
          ) : (
            <div className="bg-bg-card rounded-xl overflow-hidden divide-y divide-bg-dark">
              {songRules.map((rule) => (
                <div key={rule.id} className="flex items-center justify-between p-4">
                  <div className="font-medium text-text-cream">{rule.name}</div>
                  <button
                    onClick={() => handleDelete(rule)}
                    disabled={deletingId === rule.id}
                    className="p-2 text-text-muted hover:text-red-400 transition-colors disabled:opacity-50"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Tip */}
        <div className="mt-8 p-4 bg-bg-card/50 rounded-xl text-sm text-text-muted">
          <span className="mr-2">💡</span>
          Tip: You can block artists and songs directly while browsing your library or viewing playlists.
        </div>
      </div>
    </main>
  );
}
