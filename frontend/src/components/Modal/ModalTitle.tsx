import { cn } from '@/lib/utils';

interface ModalTitleProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalTitle = ({ children, className }: ModalTitleProps) => {
  return (
    <h3 id="modal-title" className={cn('text-lg font-bold text-foreground leading-6', className)}>
      {children}
    </h3>
  );
};
