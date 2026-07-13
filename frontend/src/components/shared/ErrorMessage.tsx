import * as React from 'react';
import { cn } from '@/lib/utils';

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
      className={cn(
        'rounded-lg border border-destructive/20 bg-destructive/10 font-medium text-destructive',
        variant === 'sm' ? 'p-3 text-sm' : 'p-2 text-xs',
        className,
      )}
      {...props}
    >
      {content}
    </div>
  );
};
