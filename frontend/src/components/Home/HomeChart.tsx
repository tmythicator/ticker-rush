import { useQuery } from '@tanstack/react-query';
import { useChart } from '@/hooks/useChart';
import { ColorType, CrosshairMode, type ChartOptions, type DeepPartial } from 'lightweight-charts';
import { useRef, useMemo } from 'react';
import type { TradeSymbol } from '@/types';
import { getHistory } from '@/lib/api';
import { queryKeys } from '@/lib/queryKeys';
import { queryConfig } from '@/lib/queryConfig';
import styles from './HomeChart.module.css';

export const HomeChart = ({ symbol }: { symbol: TradeSymbol }) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);

  const { data: history } = useQuery({
    queryKey: queryKeys.quotes.history(symbol),
    queryFn: () => getHistory({ symbol, limit: 100 }),
    enabled: !!symbol,
    ...queryConfig.homeChart,
  });

  const latestPrice = history && history.length > 0 ? history[history.length - 1].price : null;

  const chartOptions = useMemo<DeepPartial<ChartOptions>>(
    () => ({
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
      },
      grid: {
        vertLines: { visible: false },
        horzLines: { visible: false },
      },
      height: 400,
      timeScale: {
        visible: false,
        borderVisible: false,
      },
      rightPriceScale: {
        visible: false,
        borderVisible: false,
      },
      crosshair: {
        mode: CrosshairMode.Magnet,
        vertLine: {
          visible: true,
          labelVisible: false,
          style: 0,
          width: 1,
          color: 'rgba(255, 255, 255, 0.1)',
        },
        horzLine: {
          visible: false,
          labelVisible: false,
        },
      },
      handleScroll: false,
      handleScale: false,
    }),
    [],
  );

  useChart({
    chartContainerRef,
    symbol,
    quote: null,
    options: chartOptions,
  });

  return (
    <div className={styles.wrapper}>
      <div className={styles.chartLegend}>
        <div className={styles.titleRow}>
          <h3 className={styles.title}>Bitcoin</h3>
          <span className={styles.symbol}>BTC</span>
        </div>
        {latestPrice !== null && latestPrice !== undefined && (
          <div className={styles.price}>
            $
            {latestPrice.toLocaleString(undefined, {
              minimumFractionDigits: 2,
              maximumFractionDigits: 2,
            })}
          </div>
        )}
      </div>

      <div className={styles.overlayGradient} />
      <div ref={chartContainerRef} className={styles.chartContainer} />
    </div>
  );
};
