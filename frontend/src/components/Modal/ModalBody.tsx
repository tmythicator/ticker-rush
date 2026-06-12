import { cn } from '@/lib/utils';

interface ModalBodyProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalBody = ({ children, className }: ModalBodyProps) => {
  return <div className={cn('mt-2', className)}>{children}</div>;
};
