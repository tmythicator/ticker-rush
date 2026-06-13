import { createContext, useContext } from 'react';

export const ModalContext = createContext<{ onClose: () => void } | null>(null);

export const useModalContext = () => {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error('Modal compound components must be used within <Modal />');
  }
  return context;
};
