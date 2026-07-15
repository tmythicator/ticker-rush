import * as React from 'react';
import styles from './Button.module.css';
import { cva, type VariantProps } from 'class-variance-authority';

const buttonVariants = cva(styles.button, {
  variants: {
    variant: {
      default: styles.variantDefault,
      destructive: styles.variantDestructive,
      outline: styles.variantOutline,
      secondary: styles.variantSecondary,
      ghost: styles.variantGhost,
      link: styles.variantLink,
      success: styles.variantSuccess,
      ghostDestructive: styles.variantGhostDestructive,
      unstyled: styles.variantUnstyled,
    },
    size: {
      default: styles.sizeDefault,
      sm: styles.sizeSm,
      lg: styles.sizeLg,
      icon: styles.sizeIcon,
      unstyled: styles.sizeUnstyled,
    },
  },
  defaultVariants: {
    variant: 'default',
    size: 'default',
  },
});

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return <button ref={ref} className={buttonVariants({ variant, size, className })} {...props} />;
  },
);
