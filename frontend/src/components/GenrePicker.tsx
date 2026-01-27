'use client';

import { useState } from 'react';
import { ParentGenreCount, SubGenreCount } from '@/lib/api';

interface GenrePickerProps {
    parentGenres: ParentGenreCount[];
    selectedGenres: Set<string>;
    onSelectionChange: (selected: Set<string>) => void;
}

export function GenrePicker({ parentGenres, selectedGenres, onSelectionChange }: GenrePickerProps) {
    const [expandedParent, setExpandedParent] = useState<string | null>(null);

    const toggleParent = (parentName: string) => {
        setExpandedParent(expandedParent === parentName ? null : parentName);
    };

    const isParentFullySelected = (parent: ParentGenreCount) => {
        return parent.sub_genres.every(sg => selectedGenres.has(sg.name));
    };

    const isParentPartiallySelected = (parent: ParentGenreCount) => {
        const selected = parent.sub_genres.filter(sg => selectedGenres.has(sg.name));
        return selected.length > 0 && selected.length < parent.sub_genres.length;
    };

    const toggleParentSelection = (parent: ParentGenreCount) => {
        const newSelected = new Set(selectedGenres);
        if (isParentFullySelected(parent)) {
            // Deselect all sub-genres
            parent.sub_genres.forEach(sg => newSelected.delete(sg.name));
        } else {
            // Select all sub-genres
            parent.sub_genres.forEach(sg => newSelected.add(sg.name));
        }
        onSelectionChange(newSelected);
    };

    const toggleSubGenre = (subGenreName: string) => {
        const newSelected = new Set(selectedGenres);
        if (newSelected.has(subGenreName)) {
            newSelected.delete(subGenreName);
        } else {
            newSelected.add(subGenreName);
        }
        onSelectionChange(newSelected);
    };

    return (
        <div className="space-y-2">
            {parentGenres.map((parent) => (
                <div key={parent.name} className="bg-bg-card rounded-xl overflow-hidden">
                    {/* Parent Genre Row */}
                    <div
                        className="flex items-center justify-between p-4 cursor-pointer hover:bg-bg-dark/50 transition-colors"
                        onClick={() => toggleParent(parent.name)}
                    >
                        <div className="flex items-center gap-3">
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    toggleParentSelection(parent);
                                }}
                                className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${isParentFullySelected(parent)
                                        ? 'bg-accent-orange border-accent-orange'
                                        : isParentPartiallySelected(parent)
                                            ? 'border-accent-orange bg-accent-orange/30'
                                            : 'border-text-muted hover:border-text-cream'
                                    }`}
                            >
                                {(isParentFullySelected(parent) || isParentPartiallySelected(parent)) && (
                                    <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                    </svg>
                                )}
                            </button>
                            <span className="font-display text-lg text-text-cream">{parent.name}</span>
                        </div>
                        <div className="flex items-center gap-3">
                            <span className="text-sm text-text-muted">{parent.count} songs</span>
                            <svg
                                className={`w-5 h-5 text-text-muted transition-transform ${expandedParent === parent.name ? 'rotate-180' : ''
                                    }`}
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                            </svg>
                        </div>
                    </div>

                    {/* Sub-genres (expanded) */}
                    {expandedParent === parent.name && (
                        <div className="px-4 pb-4 pt-2 border-t border-bg-dark">
                            <div className="flex flex-wrap gap-2">
                                {parent.sub_genres.map((subGenre, index) => (
                                    <button
                                        key={subGenre.name}
                                        onClick={() => toggleSubGenre(subGenre.name)}
                                        className={`px-3 py-1.5 rounded-full text-sm transition-all animate-drop-in ${selectedGenres.has(subGenre.name)
                                                ? 'bg-accent-orange text-white'
                                                : 'bg-bg-dark text-text-cream hover:bg-bg-dark/70'
                                            }`}
                                        style={{ animationDelay: `${index * 50}ms` }}
                                    >
                                        {subGenre.name}
                                        <span className="ml-1.5 opacity-70">({subGenre.count})</span>
                                    </button>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
}
