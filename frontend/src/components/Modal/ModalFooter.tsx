import { cn } from '@/lib/utils';

interface ModalFooterProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalFooter = ({ children, className }: ModalFooterProps) => {
  return (
    <div className={cn('flex gap-3 pt-4 mt-6 border-t border-border/50', className)}>
      {children}
    </div>
  );
};
