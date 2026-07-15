import clsx from 'clsx';
import styles from './Modal.module.css';
import { useModalContext } from './ModalContext';

interface ModalTitleProps {
  children: React.ReactNode;
  className?: string;
  id?: string;
}

export const ModalTitle = ({ children, className, id }: ModalTitleProps) => {
  const context = useModalContext();
  const elementId = id || context?.labelId || '';

  return (
    <h3 id={elementId} className={clsx(styles.title, className)}>
      {children}
    </h3>
  );
};
