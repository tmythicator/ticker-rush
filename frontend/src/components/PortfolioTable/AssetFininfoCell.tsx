import * as React from 'react';
import styles from './AssetFininfoCell.module.css';
import { cva, type VariantProps } from 'class-variance-authority';
import clsx from 'clsx';

const cellVariants = cva(styles.cell, {
  variants: {
    variant: {
      default: styles.variantDefault,
      muted: styles.variantMuted,
      medium: styles.variantMedium,
      bold: styles.variantBold,
    },
    align: {
      center: styles.alignCenter,
      right: styles.alignRight,
    },
    trend: {
      up: styles.trendUp,
      down: styles.trendDown,
      neutral: '',
    },
  },
  defaultVariants: {
    variant: 'default',
    align: 'center',
    trend: 'neutral',
  },
});

export interface AssetFininfoCellProps
  extends
    Omit<React.TdHTMLAttributes<HTMLTableCellElement>, 'align'>,
    VariantProps<typeof cellVariants> {}

export const AssetFininfoCell = React.forwardRef<HTMLTableCellElement, AssetFininfoCellProps>(
  ({ className, variant, trend, align, ...props }, ref) => {
    return (
      <td
        ref={ref}
        className={clsx(cellVariants({ variant, trend, align }), className)}
        {...props}
      />
    );
  },
);
