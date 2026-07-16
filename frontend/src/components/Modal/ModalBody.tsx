import clsx from 'clsx';
import styles from './Modal.module.css';
interface ModalBodyProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalBody = ({ children, className }: ModalBodyProps) => {
  return <div className={clsx(styles.body, className)}>{children}</div>;
};
