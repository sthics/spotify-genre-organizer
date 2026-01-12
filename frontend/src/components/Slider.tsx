interface SliderProps {
  value: number;
  onChange: (value: number) => void;
  min: number;
  max: number;
  label?: string;
}

export function Slider({ value, onChange, min, max, label }: SliderProps) {
  const percentage = ((value - min) / (max - min)) * 100;

  return (
    <div className="w-full">
      {label && (
        <label className="block font-body text-text-muted mb-2">{label}</label>
      )}
      <div className="relative">
        <input
          type="range"
          min={min}
          max={max}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
          className="w-full h-2 bg-bg-card rounded-lg appearance-none cursor-pointer accent-accent-orange"
          style={{
            background: `linear-gradient(to right, #e85d04 0%, #e85d04 ${percentage}%, #252525 ${percentage}%, #252525 100%)`,
          }}
        />
        <div className="flex justify-between text-sm text-text-muted mt-1">
          <span>{min}</span>
          <span className="text-accent-orange font-medium text-lg">{value}</span>
          <span>{max}</span>
        </div>
      </div>
    </div>
  );
}
