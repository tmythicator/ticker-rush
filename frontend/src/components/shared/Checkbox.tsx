import * as React from 'react';
import styles from './Checkbox.module.css';

export interface CheckboxProps extends React.ComponentProps<'input'> {
  ref?: React.Ref<HTMLInputElement>;
}

export const Checkbox = ({ className, ref, ...props }: CheckboxProps) => {
  return (
    <input
      type="checkbox"
      className={`${styles.checkbox} ${className || ''}`}
      ref={ref}
      {...props}
    />
  );
};
