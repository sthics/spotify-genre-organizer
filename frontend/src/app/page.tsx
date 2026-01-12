'use client';

import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { useEffect, useState } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Home() {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setIsVisible(true);
  }, []);

  const handleConnect = () => {
    window.location.href = `${API_URL}/api/auth/login`;
  };

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4">
      {/* Hero Section */}
      <div className={`text-center max-w-2xl transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}>
        {/* Logo */}
        <div className="flex items-center justify-center gap-4 mb-8">
          <h1 className="font-display text-4xl md:text-5xl lg:text-6xl text-text-cream">
            Spotify Genre
            <br />
            Organizer
          </h1>
          <VinylIcon spinning size={80} />
        </div>

        {/* Tagline */}
        <div className="mb-12 space-y-2">
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            &ldquo;2,000 liked songs.
          </p>
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            Zero organization.
          </p>
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            Sound familiar?&rdquo;
          </p>
        </div>

        {/* CTA Button */}
        <Button
          size="lg"
          onClick={handleConnect}
          className="flex items-center gap-2 mx-auto"
        >
          <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z"/>
          </svg>
          Connect with Spotify
        </Button>
      </div>

      {/* Value Props */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-20 max-w-4xl w-full px-4">
        {[
          { title: 'Analyze', desc: 'your library' },
          { title: 'Organize', desc: 'into playlists' },
          { title: 'Enjoy', desc: 'your music' },
        ].map((item, index) => (
          <div
            key={item.title}
            className={`bg-bg-card p-6 rounded-xl text-center transition-all duration-700 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
            style={{ transitionDelay: `${800 + index * 200}ms` }}
          >
            <h3 className="font-display text-xl text-accent-orange mb-2">{item.title}</h3>
            <p className="text-text-muted">{item.desc}</p>
          </div>
        ))}
      </div>
    </main>
  );
}
