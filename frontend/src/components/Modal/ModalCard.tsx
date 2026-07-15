import { useEffect, useRef } from 'react';
import styles from './Modal.module.css';

interface ModalCardProps {
  children: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
  className?: string;
}

export const ModalCard = ({ children, size = 'sm', className }: ModalCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    cardRef.current?.focus();
  }, []);

  return (
    <div
      ref={cardRef}
      tabIndex={-1}
      className={`${styles.card} ${className || ''}`}
      data-size={size}
    >
      {children}
    </div>
  );
};
