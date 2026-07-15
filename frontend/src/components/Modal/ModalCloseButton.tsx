import { IconX } from '@/components/icons/CustomIcons';
import { Button } from '@/components/shared/Button';
import { useModalContext } from './ModalContext';
import styles from './Modal.module.css';

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
      className={`${styles.closeButton} ${className || ''}`}
    >
      <IconX className={styles.closeIcon} />
      <span className={styles.srOnly}>Close</span>
    </Button>
  );
};
