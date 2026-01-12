import { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

export function Button({
  children,
  variant = 'primary',
  size = 'md',
  className = '',
  ...props
}: ButtonProps) {
  const baseStyles = 'font-body font-medium rounded-lg transition-all duration-200 transform active:scale-95';

  const variants = {
    primary: 'bg-accent-orange hover:bg-accent-orange-hover text-white shadow-lg hover:shadow-xl',
    secondary: 'bg-bg-card hover:bg-bg-card/80 text-text-cream border border-text-muted/20',
    ghost: 'bg-transparent hover:bg-bg-card text-text-muted hover:text-text-cream',
  };

  const sizes = {
    sm: 'px-4 py-2 text-sm',
    md: 'px-6 py-3 text-base',
    lg: 'px-8 py-4 text-lg',
  };

  return (
    <button
      className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
