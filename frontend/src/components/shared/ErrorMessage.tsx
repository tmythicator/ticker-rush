import * as React from 'react';
import styles from './ErrorMessage.module.css';
import clsx from 'clsx';

export interface ErrorMessageProps extends React.HTMLAttributes<HTMLDivElement> {
  message?: string;
  variant?: 'sm' | 'xs';
}

export const ErrorMessage = ({
  className,
  message,
  children,
  variant = 'sm',
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
      className={clsx(styles.errorMessage, className)}
      data-variant={variant}
      {...props}
    >
      {content}
    </div>
  );
};
