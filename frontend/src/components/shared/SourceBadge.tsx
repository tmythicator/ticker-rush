import { type TickerSource } from '@/types';
import styles from './SourceBadge.module.css';

export interface SourceBadgeProps extends React.ComponentProps<'span'> {
  source: TickerSource;
  ref?: React.Ref<HTMLSpanElement>;
}

export const SourceBadge = ({ source, className, ref, ...props }: SourceBadgeProps) => {
  const isCoinGecko = source === 'CoinGecko' || source === 'CG';
  const label = isCoinGecko ? 'Source: CoinGecko' : 'Source: Finnhub';
  const displayLabel = isCoinGecko ? 'CG' : 'FH';
  const sourceName = isCoinGecko ? 'CoinGecko' : 'Finnhub';

  return (
    <span
      ref={ref}
      className={`${styles.badge} ${className || ''}`}
      data-source={sourceName}
      title={label}
      {...props}
    >
      {displayLabel}
    </span>
  );
};
