interface PlaylistCardProps {
  genre: string;
  songCount: number;
  spotifyUrl: string;
  index: number;
}

const genreColors: Record<string, string> = {
  'Rock': '#dc2626',
  'Pop': '#ec4899',
  'Hip-Hop': '#8b5cf6',
  'Electronic': '#06b6d4',
  'R&B': '#f59e0b',
  'Jazz': '#10b981',
  'Classical': '#6366f1',
  'Country': '#ca8a04',
  'Metal': '#374151',
  'Folk': '#84cc16',
  'Latin': '#f97316',
  'Blues': '#3b82f6',
  'Reggae': '#22c55e',
  'Punk': '#ef4444',
  'Indie': '#a855f7',
  'Soul': '#eab308',
  'Funk': '#d946ef',
  'World': '#14b8a6',
  'Other': '#6b7280',
};

export function PlaylistCard({ genre, songCount, spotifyUrl, index }: PlaylistCardProps) {
  const accentColor = genreColors[genre] || genreColors['Other'];

  return (
    <div
      className="bg-bg-card rounded-xl p-4 transform hover:-translate-y-1 hover:shadow-xl transition-all duration-200 animate-drop-in cursor-pointer group"
      style={{ animationDelay: `${index * 100}ms` }}
      onClick={() => window.open(spotifyUrl, '_blank')}
    >
      {/* Accent bar */}
      <div
        className="h-1 rounded-full mb-3"
        style={{ backgroundColor: accentColor }}
      />

      {/* Content */}
      <div className="flex items-center gap-2 mb-2">
        <span className="text-lg">&#9835;</span>
        <h3 className="font-display text-lg text-text-cream group-hover:text-accent-orange transition-colors">
          {genre}
        </h3>
      </div>

      <p className="text-text-muted text-sm mb-3">
        {songCount} songs
      </p>

      <button className="text-sm text-accent-orange hover:text-accent-orange-hover transition-colors">
        Open in Spotify &rarr;
      </button>
    </div>
  );
}
