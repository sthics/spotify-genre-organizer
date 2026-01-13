'use client';

import { useState, useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { useUser } from '@/hooks/useUser';
import { getSettings, updateSettings, UserSettings } from '@/lib/api';

const EXAMPLE_GENRES = ['Rock', 'Electronic', 'Jazz'];

const TOKEN_BUTTONS = [
  { token: '{genre}', label: 'Genre', required: true },
  { token: '{username}', label: 'Username', required: false },
  { token: '{date}', label: 'Date', required: false },
];

const FREE_TIER_FOOTER = '\n\nOrganized by Spotify Genre Organizer';

export default function Settings() {
  const { user, loading: userLoading } = useUser();
  const router = useRouter();

  const [settings, setSettings] = useState<UserSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [nameTemplate, setNameTemplate] = useState('{genre}');
  const [descriptionTemplate, setDescriptionTemplate] = useState('A playlist of {genre} tracks from my library.');

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const data = await getSettings();
        setSettings(data);
        setNameTemplate(data.name_template);
        setDescriptionTemplate(data.description_template);
      } catch (err) {
        console.error('Failed to fetch settings:', err);
        // Use defaults if fetch fails
      } finally {
        setLoading(false);
      }
    };

    fetchSettings();
  }, []);

  const isValid = useMemo(() => {
    return nameTemplate.includes('{genre}');
  }, [nameTemplate]);

  const previews = useMemo(() => {
    const username = user?.display_name || 'User';
    const date = new Date().toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });

    return EXAMPLE_GENRES.map((genre) => {
      let name = nameTemplate
        .replace(/{genre}/g, genre)
        .replace(/{username}/g, username)
        .replace(/{date}/g, date);

      return name;
    });
  }, [nameTemplate, user?.display_name]);

  const descriptionPreview = useMemo(() => {
    const username = user?.display_name || 'User';
    const date = new Date().toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });

    let description = descriptionTemplate
      .replace(/{genre}/g, EXAMPLE_GENRES[0])
      .replace(/{username}/g, username)
      .replace(/{date}/g, date);

    // Add free tier footer if not premium
    if (!settings?.is_premium) {
      description += FREE_TIER_FOOTER;
    }

    return description;
  }, [descriptionTemplate, user?.display_name, settings?.is_premium]);

  const insertToken = (token: string, field: 'name' | 'description') => {
    if (field === 'name') {
      setNameTemplate((prev) => prev + token);
    } else {
      setDescriptionTemplate((prev) => prev + token);
    }
  };

  const handleSave = async () => {
    if (!isValid) {
      setError('Name template must contain {genre}');
      return;
    }

    setSaving(true);
    setError(null);

    try {
      const updated = await updateSettings(nameTemplate, descriptionTemplate);
      setSettings(updated);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (err) {
      console.error('Failed to save settings:', err);
      setError('Failed to save settings. Please try again.');
    } finally {
      setSaving(false);
    }
  };

  const handleBack = () => {
    router.back();
  };

  if (userLoading || loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12">
      {/* Header */}
      <div className="w-full max-w-xl mb-8">
        <button
          onClick={handleBack}
          className="flex items-center gap-2 text-text-muted hover:text-text-cream transition-colors mb-6"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>

        <div className="flex items-center gap-3">
          <VinylIcon size={40} />
          <div>
            <h1 className="font-display text-3xl text-text-cream">Label Press</h1>
            <p className="text-text-muted">Customize how your playlists are named</p>
          </div>
        </div>
      </div>

      {/* Main Card */}
      <div className="w-full max-w-xl bg-bg-card rounded-2xl p-8 shadow-xl">
        {/* Name Template Section */}
        <div className="mb-8">
          <label className="block font-display text-xl text-text-cream mb-2">
            Playlist Name
          </label>
          <p className="text-text-muted text-sm mb-4">
            Must include {'{genre}'} token
          </p>

          {/* Token Buttons */}
          <div className="flex flex-wrap gap-2 mb-3">
            {TOKEN_BUTTONS.map(({ token, label, required }) => (
              <button
                key={token}
                onClick={() => insertToken(token, 'name')}
                className={`px-3 py-1.5 text-sm rounded-lg border transition-colors ${
                  required
                    ? 'border-accent-orange text-accent-orange hover:bg-accent-orange hover:text-white'
                    : 'border-text-muted/30 text-text-muted hover:border-text-cream hover:text-text-cream'
                }`}
              >
                {token}
                {required && <span className="text-xs ml-1">*</span>}
              </button>
            ))}
          </div>

          {/* Name Input */}
          <input
            type="text"
            value={nameTemplate}
            onChange={(e) => setNameTemplate(e.target.value)}
            className={`w-full px-4 py-3 bg-bg-dark rounded-lg border font-mono text-text-cream placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-accent-orange transition-all ${
              !isValid ? 'border-red-500' : 'border-text-muted/20'
            }`}
            placeholder="{genre} Playlist"
          />
          {!isValid && (
            <p className="text-red-400 text-sm mt-2">
              Name template must contain {'{genre}'}
            </p>
          )}
        </div>

        {/* Preview Section */}
        <div className="mb-8">
          <label className="block font-display text-lg text-text-cream mb-3">
            Preview
          </label>
          <div className="space-y-2">
            {previews.map((preview, index) => (
              <div
                key={index}
                className="flex items-center gap-3 px-4 py-2 bg-bg-dark rounded-lg"
              >
                <VinylIcon size={24} />
                <span className="text-text-cream font-mono">{preview}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Description Template Section */}
        <div className="mb-8">
          <label className="block font-display text-xl text-text-cream mb-2">
            Playlist Description
          </label>
          <p className="text-text-muted text-sm mb-4">
            Optional tokens: {'{genre}'}, {'{username}'}, {'{date}'}
          </p>

          {/* Token Buttons */}
          <div className="flex flex-wrap gap-2 mb-3">
            {TOKEN_BUTTONS.map(({ token, label }) => (
              <button
                key={token}
                onClick={() => insertToken(token, 'description')}
                className="px-3 py-1.5 text-sm rounded-lg border border-text-muted/30 text-text-muted hover:border-text-cream hover:text-text-cream transition-colors"
              >
                {token}
              </button>
            ))}
          </div>

          {/* Description Textarea */}
          <textarea
            value={descriptionTemplate}
            onChange={(e) => setDescriptionTemplate(e.target.value)}
            rows={3}
            className="w-full px-4 py-3 bg-bg-dark rounded-lg border border-text-muted/20 font-mono text-text-cream placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-accent-orange transition-all resize-none"
            placeholder="A playlist of {genre} music..."
          />
        </div>

        {/* Description Preview */}
        <div className="mb-8">
          <label className="block font-display text-lg text-text-cream mb-3">
            Description Preview
          </label>
          <div className="px-4 py-3 bg-bg-dark rounded-lg">
            <p className="text-text-muted font-mono text-sm whitespace-pre-wrap">
              {descriptionPreview}
            </p>
          </div>
        </div>

        {/* Free Tier Notice */}
        {!settings?.is_premium && (
          <div className="mb-8 p-4 bg-bg-dark rounded-lg border border-text-muted/20">
            <div className="flex items-start gap-3">
              <svg className="w-5 h-5 text-accent-orange flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <p className="text-text-cream text-sm font-medium">Free Tier</p>
                <p className="text-text-muted text-sm mt-1">
                  The following text will be added to all playlist descriptions:
                </p>
                <p className="text-text-muted text-xs mt-2 font-mono italic">
                  &quot;{FREE_TIER_FOOTER.trim()}&quot;
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg">
            <p className="text-red-400 text-sm">{error}</p>
          </div>
        )}

        {/* Save Button */}
        <Button
          size="lg"
          className="w-full flex items-center justify-center gap-2"
          onClick={handleSave}
          disabled={saving || !isValid}
        >
          {saving ? (
            <>
              <VinylIcon spinning size={24} />
              Saving...
            </>
          ) : saved ? (
            <>
              <svg className="w-6 h-6 text-success-green" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              Saved!
            </>
          ) : (
            'Save Settings'
          )}
        </Button>
      </div>
    </main>
  );
}
