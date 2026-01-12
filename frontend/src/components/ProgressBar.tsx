interface ProgressBarProps {
  progress: number;
  label?: string;
}

export function ProgressBar({ progress, label }: ProgressBarProps) {
  return (
    <div className="w-full">
      <div className="h-2 bg-bg-card rounded-full overflow-hidden">
        <div
          className="h-full bg-text-cream transition-all duration-300 ease-out"
          style={{ width: `${Math.min(100, Math.max(0, progress))}%` }}
        />
      </div>
      {label && (
        <p className="text-center text-text-muted mt-2">{label}</p>
      )}
    </div>
  );
}
