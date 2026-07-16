import { cva, type VariantProps } from 'class-variance-authority';
import clsx from 'clsx';
import * as React from 'react';
import styles from './Label.module.css';

const labelVariants = cva(styles.label, {
  variants: {
    variant: {
      default: styles.variantDefault,
      muted: styles.variantMuted,
      error: styles.variantError,
      plain: styles.variantPlain,
    },
  },
  defaultVariants: {
    variant: 'default',
  },
});

export interface LabelProps
  extends React.LabelHTMLAttributes<HTMLLabelElement>, VariantProps<typeof labelVariants> {}

export const Label = React.forwardRef<HTMLLabelElement, LabelProps>(
  ({ className, variant, htmlFor, ...props }, ref) => {
    return (
      <label
        ref={ref}
        htmlFor={htmlFor}
        className={clsx(labelVariants({ variant }), className)}
        {...props}
      />
    );
  },
);
