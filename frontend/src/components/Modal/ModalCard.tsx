import { useEffect, useRef } from 'react';
import styles from './Modal.module.css';
import { cva, type VariantProps } from 'class-variance-authority';
import clsx from 'clsx';

const modalVariants = cva(styles.card, {
  variants: {
    size: {
      sm: styles.sizeSm,
      md: styles.sizeMd,
      lg: styles.sizeLg,
      xl: styles.sizeXl,
      '2xl': styles.size2Xl,
      full: styles.sizeFull,
    },
  },
  defaultVariants: {
    size: 'sm',
  },
});

export interface ModalCardProps
  extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof modalVariants> {}

export const ModalCard = ({ children, size, className, ...props }: ModalCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    cardRef.current?.focus();
  }, []);

  return (
    <div
      ref={cardRef}
      tabIndex={-1}
      className={clsx(modalVariants({ size }), className)}
      {...props}
    >
      {children}
    </div>
  );
};
