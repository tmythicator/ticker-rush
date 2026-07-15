import { useChart } from '@/hooks/useChart';
import { type Quote, type TradeSymbol } from '@/types';
import { useRef } from 'react';
import { ChartSymbolIndicator } from './ChartSymbolIndicator';
import { ChartSymbolPicker } from './ChartSymbolPicker';
import styles from './MarketChart.module.css';

interface MarketChartProps {
  symbol: TradeSymbol | null;
  onSymbolChange: (s: TradeSymbol) => void;
  quote: Quote | null;
  isLoading: boolean;
  isError: boolean;
}

export const MarketChart = ({
  symbol,
  onSymbolChange,
  quote,
  isLoading,
  isError,
}: MarketChartProps) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);

  useChart({ chartContainerRef, quote, symbol });

  if (!symbol) {
    return (
      <div className={styles.emptyState}>
        No active ticker selected. Check configuration.
      </div>
    );
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.controlsWrapper}>
        <div className={styles.controlsBox}>
          <ChartSymbolPicker symbol={symbol} onSymbolChange={onSymbolChange} />
          <div className={styles.divider}></div>
          <ChartSymbolIndicator quote={quote} isLoading={isLoading} isError={isError} />
        </div>
      </div>
      <div className={styles.wrapper}>
        <div ref={chartContainerRef} className={styles.chartContainer} />
      </div>
    </div>
  );
};
