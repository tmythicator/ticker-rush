import { type TickerSource } from '@/types';
import styles from './SourceBadge.module.css';
import { cva, type VariantProps } from 'class-variance-authority';
import React from 'react';
import { getSourceBadgeConfig } from '@/lib/utils';

const badgeVariants = cva(styles.badge, {
  variants: {
    variant: {
      CoinGecko: styles.variantCoinGecko,
      Finnhub: styles.variantFinnhub,
    },
  },
});

export interface SourceBadgeProps
  extends React.ComponentProps<'span'>, VariantProps<typeof badgeVariants> {
  source: TickerSource;
}

export const SourceBadge = React.forwardRef<HTMLSpanElement, SourceBadgeProps>(
  ({ source, className, ...props }, ref) => {
    const { variant, label, title } = getSourceBadgeConfig(source);
    return (
      <span ref={ref} className={badgeVariants({ variant, className })} title={title} {...props}>
        {label}
      </span>
    );
  },
);
