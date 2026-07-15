import * as React from 'react';
import styles from './Input.module.css';

export interface InputProps extends Omit<React.ComponentProps<'input'>, 'size'> {
  variant?: 'default' | 'error' | 'unstyled';
  size?: 'default' | 'sm' | 'lg' | 'unstyled';
  ref?: React.Ref<HTMLInputElement>;
}

export const Input = ({
  className,
  type,
  variant = 'default',
  size = 'default',
  ref,
  ...props
}: InputProps) => {
  const combinedClassName = className ? `${styles.input} ${className}` : styles.input;

  return (
    <input
      type={type}
      className={combinedClassName}
      data-variant={variant}
      data-size={size}
      ref={ref}
      {...props}
    />
  );
};
