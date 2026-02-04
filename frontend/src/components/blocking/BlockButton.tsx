'use client';

import { useState } from 'react';

interface BlockButtonProps {
  type: 'artist' | 'song';
  spotifyId: string;
  name: string;
  imageUrl?: string;
  onBlock: () => void;
}

export function BlockButton({ onBlock }: BlockButtonProps) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <button
      onClick={(e) => {
        e.stopPropagation();
        onBlock();
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className={`p-1.5 rounded-full transition-all duration-200 ${
        isHovered
          ? 'bg-red-500/20 text-red-400'
          : 'text-text-muted hover:text-red-400'
      }`}
      title="Block"
    >
      <svg
        className="w-4 h-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
        />
      </svg>
    </button>
  );
}
