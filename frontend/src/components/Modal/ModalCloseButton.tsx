import { IconX } from '@/components/icons/CustomIcons';
import { Button } from '@/components/shared/Button';
import { useModalContext } from './ModalContext';
import { cn } from '@/lib/utils';

interface ModalCloseButtonProps {
  className?: string;
}

export const ModalCloseButton = ({ className }: ModalCloseButtonProps) => {
  const { onClose } = useModalContext();

  return (
    <Button
      onClick={onClose}
      variant="ghost"
      size="icon"
      className={cn('h-8 w-8 shrink-0 rounded-full', className)}
    >
      <IconX className="h-5 w-5" />
      <span className="sr-only">Close</span>
    </Button>
  );
};
