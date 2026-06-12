import { cn } from '@/lib/utils';

interface ModalHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalHeader = ({ children, className }: ModalHeaderProps) => {
  return (
    <div className={cn('flex items-center justify-between mb-5 gap-4', className)}>{children}</div>
  );
};
