import styles from './Modal.module.css';

interface ModalBodyProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalBody = ({ children, className }: ModalBodyProps) => {
  return <div className={`${styles.body} ${className || ''}`}>{children}</div>;
};
