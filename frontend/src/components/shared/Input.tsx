import * as React from 'react';
import styles from './Input.module.css';
import { cva, type VariantProps } from 'class-variance-authority';

const inputVariants = cva(styles.input, {
  variants: {
    variant: {
      default: styles.variantDefault,
      error: styles.variantError,
      unstyled: styles.variantUnstyled,
    },
    size: {
      default: styles.sizeDefault,
      sm: styles.sizeSm,
      lg: styles.sizeLg,
      unstyled: styles.sizeUnstyled,
    },
  },
  defaultVariants: {
    variant: 'default',
    size: 'default',
  },
});

export interface InputProps
  extends Omit<React.ComponentProps<'input'>, 'size'>, VariantProps<typeof inputVariants> {}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, variant, size, ...props }: InputProps, ref) => {
    return <input className={inputVariants({ variant, size, className })} ref={ref} {...props} />;
  },
);
