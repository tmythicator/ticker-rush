import styles from './Modal.module.css';

interface ModalTitleProps {
  children: React.ReactNode;
  className?: string;
}

export const ModalTitle = ({ children, className }: ModalTitleProps) => {
  return (
    <h3 id="modal-title" className={`${styles.title} ${className || ''}`}>
      {children}
    </h3>
  );
};
