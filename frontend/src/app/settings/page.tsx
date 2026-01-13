'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/Button';
import { getSettings, updateSettings, UserSettings } from '@/lib/api';

export default function Settings() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [, setSettings] = useState<UserSettings | null>(null);

  // Local state for inputs
  const [nameTemplate, setNameTemplate] = useState('');
  const [descTemplate, setDescTemplate] = useState('');

  useEffect(() => {
    async function load() {
      try {
        const data = await getSettings();
        setSettings(data);
        setNameTemplate(data.name_template);
        setDescTemplate(data.description_template);
      } catch (err) {
        console.error('Failed to load settings', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const handleSave = async () => {
    if (!nameTemplate.includes('{genre}')) {
      alert('Playlist name must contain {genre} placeholder');
      return;
    }

    setSaving(true);
    try {
      const updated = await updateSettings(nameTemplate, descTemplate);
      setSettings(updated);
      alert('Settings saved!');
    } catch (err) {
      console.error(err);
      alert('Failed to save settings');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <div className="text-text-cream animate-pulse">Loading settings...</div>
      </main>
    );
  }

  // Preview generation
  const previewName = nameTemplate
    .replace('{genre}', 'Rock')
    .replace('{year}', new Date().getFullYear().toString());

  const previewDesc = descTemplate
    .replace('{genre}', 'Rock')
    .replace('{year}', new Date().getFullYear().toString());

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12">
      <div className="w-full max-w-2xl">
        <button
          onClick={() => router.push('/dashboard')}
          className="text-text-muted hover:text-text-cream transition-colors mb-8"
        >
          &larr; Back to Dashboard
        </button>

        <h1 className="font-display text-3xl text-text-cream mb-8">
          Label Settings
        </h1>

        <div className="bg-bg-card rounded-xl p-8 shadow-xl space-y-8">
          {/* Naming Template */}
          <div>
            <label className="block text-text-cream font-medium mb-2">
              Playlist Name Pattern
            </label>
            <p className="text-sm text-text-muted mb-4">
              How should we name your new crates?
            </p>
            <input
              type="text"
              value={nameTemplate}
              onChange={(e) => setNameTemplate(e.target.value)}
              className="w-full bg-bg-dark border border-gray-700 rounded-lg px-4 py-3 text-text-cream focus:ring-2 focus:ring-accent-orange outline-none transition-all"
              placeholder="{genre} by Organizer"
            />
            <div className="mt-2 flex gap-2 text-xs">
              <span
                className="bg-gray-700 text-gray-300 px-2 py-1 rounded cursor-pointer hover:bg-gray-600 transition-colors"
                onClick={() => setNameTemplate(prev => prev + ' {genre}')}
              >
                + &#123;genre&#125;
              </span>
              <span
                className="bg-gray-700 text-gray-300 px-2 py-1 rounded cursor-pointer hover:bg-gray-600 transition-colors"
                onClick={() => setNameTemplate(prev => prev + ' {year}')}
              >
                + &#123;year&#125;
              </span>
            </div>
          </div>

          {/* Description Template */}
          <div>
            <label className="block text-text-cream font-medium mb-2">
              Description Pattern
            </label>
            <textarea
              value={descTemplate}
              onChange={(e) => setDescTemplate(e.target.value)}
              className="w-full bg-bg-dark border border-gray-700 rounded-lg px-4 py-3 text-text-cream focus:ring-2 focus:ring-accent-orange outline-none transition-all"
              rows={3}
              placeholder="Organized by Spotify Genre Organizer"
            />
            <div className="mt-2 flex gap-2 text-xs">
              <span
                className="bg-gray-700 text-gray-300 px-2 py-1 rounded cursor-pointer hover:bg-gray-600 transition-colors"
                onClick={() => setDescTemplate(prev => prev + ' {genre}')}
              >
                + &#123;genre&#125;
              </span>
              <span
                className="bg-gray-700 text-gray-300 px-2 py-1 rounded cursor-pointer hover:bg-gray-600 transition-colors"
                onClick={() => setDescTemplate(prev => prev + ' {year}')}
              >
                + &#123;year&#125;
              </span>
            </div>
          </div>

          {/* Preview Card */}
          <div className="bg-bg-dark rounded-lg p-6 border border-gray-700">
            <h3 className="text-xs font-bold text-text-muted uppercase tracking-wider mb-4">
              Preview
            </h3>
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 bg-gradient-to-br from-orange-500 to-red-600 rounded shadow-lg flex-shrink-0" />
              <div>
                <h4 className="font-display text-lg text-text-cream">
                  {previewName}
                </h4>
                <p className="text-sm text-text-muted">
                  {previewDesc}
                </p>
              </div>
            </div>
          </div>

          <div className="pt-4 flex justify-end">
            <Button onClick={handleSave} disabled={saving}>
              {saving ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </div>
      </div>
    </main>
  );
}
