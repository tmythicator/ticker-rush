import * as React from 'react';
import styles from './ErrorMessage.module.css';
import { cva, type VariantProps } from 'class-variance-authority';

const errorVariants = cva(styles.errorMessage, {
  variants: {
    variant: {
      sm: styles.sizeSm,
      xs: styles.sizeXs,
    },
  },
  defaultVariants: {
    variant: 'sm',
  },
});

export interface ErrorMessageProps
  extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof errorVariants> {
  message?: string;
}

export const ErrorMessage = ({
  className,
  message,
  children,
  variant,
  ...props
}: ErrorMessageProps) => {
  let content = message || children;
  if (!content) return null;

  if (typeof content === 'string' && content.length > 0) {
    content = content.charAt(0).toUpperCase() + content.slice(1);
  }

  return (
    <div
      role="alert"
      data-testid="error-message"
      className={errorVariants({ variant, className })}
      {...props}
    >
      {content}
    </div>
  );
};
