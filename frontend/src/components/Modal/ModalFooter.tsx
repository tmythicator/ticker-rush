import styles from './Modal.module.css';

interface ModalFooterProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalFooter = ({ children, className }: ModalFooterProps) => {
  return <div className={`${styles.footer} ${className || ''}`}>{children}</div>;
};
