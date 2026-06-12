import { useEffect, useRef } from 'react';
import { cn } from '@/lib/utils';

interface ModalCardProps {
  children: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
  className?: string;
}

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  full: 'max-w-full m-4 h-[calc(100vh-2rem)]',
};

export const ModalCard = ({ children, size = 'sm', className }: ModalCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    cardRef.current?.focus();
  }, []);

  return (
    <div
      ref={cardRef}
      tabIndex={-1}
      className={cn(
        'relative w-full transform overflow-hidden rounded-lg bg-card p-6 text-left shadow-xl transition-all border border-border animate-in fade-in zoom-in-95 duration-200 outline-none z-10',
        sizeClasses[size],
        className,
      )}
    >
      {children}
    </div>
  );
};
