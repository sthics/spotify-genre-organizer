interface VinylIconProps {
  className?: string;
  spinning?: boolean;
  size?: number;
}

export function VinylIcon({ className = "", spinning = false, size = 64 }: VinylIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 100 100"
      className={`${spinning ? 'animate-spin-slow' : ''} ${className}`}
    >
      {/* Outer ring */}
      <circle cx="50" cy="50" r="48" fill="#1a1a1a" stroke="#333" strokeWidth="2" />

      {/* Grooves */}
      <circle cx="50" cy="50" r="40" fill="none" stroke="#2a2a2a" strokeWidth="1" />
      <circle cx="50" cy="50" r="35" fill="none" stroke="#252525" strokeWidth="1" />
      <circle cx="50" cy="50" r="30" fill="none" stroke="#2a2a2a" strokeWidth="1" />
      <circle cx="50" cy="50" r="25" fill="none" stroke="#252525" strokeWidth="1" />
      <circle cx="50" cy="50" r="20" fill="none" stroke="#2a2a2a" strokeWidth="1" />

      {/* Label */}
      <circle cx="50" cy="50" r="15" fill="#e85d04" />
      <circle cx="50" cy="50" r="12" fill="#ff6b0a" />

      {/* Center hole */}
      <circle cx="50" cy="50" r="3" fill="#1a1a1a" />

      {/* Shine effect */}
      <ellipse cx="35" cy="35" rx="8" ry="4" fill="rgba(255,255,255,0.05)" transform="rotate(-45 35 35)" />
    </svg>
  );
}
