import * as React from 'react';
import styles from './Label.module.css';

export interface LabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
  variant?: 'default' | 'muted' | 'error';
  ref?: React.Ref<HTMLLabelElement>;
}

export const Label = ({
  className,
  variant = 'default',
  ref,
  ...props
}: LabelProps) => {
  const combinedClassName = className ? `${styles.label} ${className}` : styles.label;

  return (
    <label
      ref={ref}
      className={combinedClassName}
      data-variant={variant}
      {...props}
    />
  );
};
