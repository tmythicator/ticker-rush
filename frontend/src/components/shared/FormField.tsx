import * as React from 'react';
import { Label } from '@/components/shared/Label';
import styles from './FormField.module.css';
import clsx from 'clsx';

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
    <div ref={ref} className={clsx(styles.formField, className)} {...props}>
      {label && <Label htmlFor={htmlFor}>{label}</Label>}
      {children}
      {error && (
        <p data-testid="field-error" className={styles.error}>
          {error}
        </p>
      )}
    </div>
  );
};
