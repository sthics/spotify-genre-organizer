'use client';

import { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { Toast, ToastData } from '@/components/Toast';

interface ToastContextType {
  showToast: (message: string, type: 'success' | 'error', action?: { label: string; onClick: () => void }) => void;
}

const ToastContext = createContext<ToastContextType | null>(null);

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const showToast = useCallback((
    message: string,
    type: 'success' | 'error',
    action?: { label: string; onClick: () => void }
  ) => {
    const id = crypto.randomUUID();
    setToasts(prev => [...prev, { id, message, type, action }]);
  }, []);

  const dismissToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(t => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
        {toasts.map(toast => (
          <Toast key={toast.id} toast={toast} onDismiss={dismissToast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}
