interface GenreTagProps {
  genre: string;
  index: number;
}

export function GenreTag({ genre, index }: GenreTagProps) {
  return (
    <span
      className="inline-block px-3 py-1 bg-bg-dark text-text-cream rounded-full text-sm font-body animate-bounce-in"
      style={{ animationDelay: `${index * 100}ms` }}
    >
      {genre}
    </span>
  );
}
