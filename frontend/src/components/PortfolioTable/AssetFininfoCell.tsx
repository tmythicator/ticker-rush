import * as React from 'react';
import styles from './AssetFininfoCell.module.css';
import { cva, type VariantProps } from 'class-variance-authority';

const cellVariants = cva(styles.cell, {
  variants: {
    variant: {
      default: styles.variantDefault,
      muted: styles.variantMuted,
      medium: styles.variantMedium,
      bold: styles.variantBold,
    },
    trend: {
      up: styles.trendUp,
      down: styles.trendDown,
      neutral: '',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
});

export interface AssetFininfoCellProps
  extends React.TdHTMLAttributes<HTMLTableCellElement>, VariantProps<typeof cellVariants> {}

export const AssetFininfoCell = React.forwardRef<HTMLTableCellElement, AssetFininfoCellProps>(
  ({ className, variant, trend, ...props }, ref) => {
    return <td ref={ref} className={cellVariants({ variant, trend, className })} {...props} />;
  },
);
