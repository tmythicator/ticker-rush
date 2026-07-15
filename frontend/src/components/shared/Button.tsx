import * as React from 'react';
import styles from './Button.module.css';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?:
    | 'default'
    | 'destructive'
    | 'outline'
    | 'secondary'
    | 'ghost'
    | 'link'
    | 'success'
    | 'ghostDestructive'
    | 'unstyled';
  size?: 'default' | 'sm' | 'lg' | 'icon' | 'unstyled';
  ref?: React.Ref<HTMLButtonElement>;
}

export const Button = ({
  className,
  variant = 'default',
  size = 'default',
  ref,
  ...props
}: ButtonProps) => {
  return (
    <button
      className={`${styles.button} ${className || ''}`}
      data-variant={variant}
      data-size={size}
      ref={ref}
      {...props}
    />
  );
};
