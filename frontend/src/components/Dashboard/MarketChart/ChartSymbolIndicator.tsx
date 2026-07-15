import { usePriceColor } from '@/hooks/usePriceColor';
import { type Quote } from '@/types';
import styles from './MarketChart.module.css';

interface ChartSymbolIndicatorProps {
  quote: Quote | null;
  isLoading: boolean;
  isError: boolean;
}

export const ChartSymbolIndicator = ({ quote, isLoading, isError }: ChartSymbolIndicatorProps) => {
  const price = quote?.price;
  const isClosed = quote?.is_closed;
  const priceColorStatus = usePriceColor(price);

  const priceColorClass =
    priceColorStatus === 'up'
      ? styles.priceUp
      : priceColorStatus === 'down'
        ? styles.priceDown
        : styles.priceNeutral;

  return (
    <div className={styles.indicator}>
      {isLoading ? (
        <div className={styles.pulseLoader}></div>
      ) : isError ? (
        <span className={styles.offlineTag}>OFFLINE</span>
      ) : (
        <>
          <span className={`${styles.price} ${priceColorClass}`}>
            {price ? `$${price.toFixed(2)}` : '—'}
          </span>
          {isClosed ? (
            <span className={styles.statusLabel} data-status="closed">
              Market Closed
            </span>
          ) : (
            <span className={styles.statusLabel} data-status="live">
              Live
            </span>
          )}
        </>
      )}
    </div>
  );
};
