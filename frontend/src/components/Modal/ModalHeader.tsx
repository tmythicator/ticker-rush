import clsx from 'clsx';
import styles from './Modal.module.css';

interface ModalHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalHeader = ({ children, className }: ModalHeaderProps) => {
  return <div className={clsx(styles.header, className)}>{children}</div>;
};
