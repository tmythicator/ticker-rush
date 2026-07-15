import React, { useEffect, useId, useRef } from 'react';
import { createPortal } from 'react-dom';
import { ModalContext } from './ModalContext';
import styles from './Modal.module.css';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

let activeModalsCount = 0;

export const Modal = ({ isOpen, onClose, children }: ModalProps) => {
  const modalRef = useRef<HTMLDivElement>(null);
  const previouslyFocusedElement = useRef<HTMLElement | null>(null);
  const modalLabelId = useId();

  useEffect(() => {
    if (!isOpen) return;

    activeModalsCount++;
    document.body.style.overflow = 'hidden';
    previouslyFocusedElement.current = document.activeElement as HTMLElement;

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };

    if (modalRef.current) {
      modalRef.current.focus();
    }
    document.addEventListener('keydown', handleEscape);

    return () => {
      document.removeEventListener('keydown', handleEscape);

      activeModalsCount--;
      if (activeModalsCount === 0) {
        document.body.style.overflow = 'unset';
      }

      if (previouslyFocusedElement.current) {
        previouslyFocusedElement.current.focus();
      }
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return createPortal(
    <ModalContext.Provider value={{ onClose }}>
      <div
        ref={modalRef}
        className={styles.modalContainer}
        role="dialog"
        aria-modal="true"
        aria-labelledby={modalLabelId}
        tabIndex={-1}
      >
        <div className={styles.modalOverlay} onClick={onClose} aria-hidden="true" />
        {children}
      </div>
    </ModalContext.Provider>,
    document.body,
  );
};
