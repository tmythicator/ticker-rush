import * as React from 'react';
import { Label } from '@/components/shared/Label';
import { cn } from '@/lib/utils';

export interface FormFieldProps extends React.HTMLAttributes<HTMLDivElement> {
  label?: string;
  htmlFor?: string;
  error?: string;
  ref?: React.Ref<HTMLDivElement>;
}

export const FormField = ({
  label,
  htmlFor,
  error,
  className,
  children,
  ref,
  ...props
}: FormFieldProps) => {
  return (
    <div ref={ref} className={cn('space-y-2', className)} {...props}>
      {label && <Label htmlFor={htmlFor}>{label}</Label>}
      {children}
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
};
