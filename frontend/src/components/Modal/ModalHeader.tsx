import { cn } from '@/lib/utils';

interface ModalHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalHeader = ({ children, className }: ModalHeaderProps) => {
  return (
    <div className={cn('mb-5 flex items-center justify-between gap-4', className)}>{children}</div>
  );
};
