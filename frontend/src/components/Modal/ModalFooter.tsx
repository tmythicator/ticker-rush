import { cn } from '@/lib/utils';

interface ModalFooterProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalFooter = ({ children, className }: ModalFooterProps) => {
  return (
    <div className={cn('mt-6 flex gap-3 border-t border-border/50 pt-4', className)}>
      {children}
    </div>
  );
};
