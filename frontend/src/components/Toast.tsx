'use client';

import { useEffect, useState } from 'react';

export interface ToastData {
  id: string;
  message: string;
  type: 'success' | 'error';
  action?: {
    label: string;
    onClick: () => void;
  };
}

interface ToastProps {
  toast: ToastData;
  onDismiss: (id: string) => void;
}

export function Toast({ toast, onDismiss }: ToastProps) {
  const [isExiting, setIsExiting] = useState(false);

  useEffect(() => {
    if (toast.type === 'success') {
      const timer = setTimeout(() => {
        setIsExiting(true);
        setTimeout(() => onDismiss(toast.id), 300);
      }, 4000);
      return () => clearTimeout(timer);
    }
  }, [toast.id, toast.type, onDismiss]);

  const handleDismiss = () => {
    setIsExiting(true);
    setTimeout(() => onDismiss(toast.id), 300);
  };

  const borderColor = toast.type === 'success' ? 'border-l-success-green' : 'border-l-red-500';

  return (
    <div
      className={`
        bg-bg-card rounded-lg shadow-xl border-l-4 ${borderColor}
        p-4 min-w-[280px] max-w-[380px]
        transform transition-all duration-300
        ${isExiting ? 'opacity-0 translate-x-4' : 'opacity-100 translate-x-0'}
      `}
    >
      <div className="flex items-start justify-between gap-3">
        <p className="text-text-cream text-sm">{toast.message}</p>
        <button
          onClick={handleDismiss}
          className="text-text-muted hover:text-text-cream transition-colors"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      {toast.action && (
        <button
          onClick={toast.action.onClick}
          className="mt-2 text-accent-orange hover:text-accent-orange-hover text-sm font-medium transition-colors"
        >
          {toast.action.label}
        </button>
      )}
    </div>
  );
}
