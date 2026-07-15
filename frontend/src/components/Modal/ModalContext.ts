import { createContext, useContext } from 'react';

interface ModalContextType {
  onClose: () => void;
  labelId?: string;
}

export const ModalContext = createContext<ModalContextType | null>(null);

export const useModalContext = () : Partial<ModalContextType> => {
  const context = useContext(ModalContext);
  if (!context) {
    return {};
  }
  return context;
};
