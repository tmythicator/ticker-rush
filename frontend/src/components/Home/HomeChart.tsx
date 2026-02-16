import { useQuery } from '@tanstack/react-query';
import { useChart } from '@/hooks/useChart';
import { ColorType, CrosshairMode, type ChartOptions, type DeepPartial } from 'lightweight-charts';
import { useRef, useMemo } from 'react';
import type { TradeSymbol } from '@/types';
import { getHistory } from '@/lib/api';
import { QUERY_KEY_HISTORY } from '@/lib/queryKeys';

export const HomeChart = ({ symbol }: { symbol: TradeSymbol }) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);

  const { data: history } = useQuery({
    queryKey: QUERY_KEY_HISTORY(symbol),
    queryFn: () => getHistory(symbol, 100),
    enabled: !!symbol,
    staleTime: 1000 * 60 * 3,
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
    <div className="w-full h-[400px] relative">
      {/* Overlay Info */}
      <div className="absolute top-6 left-6 z-20 pointer-events-none">
        <div className="flex items-baseline gap-2">
          <h3 className="text-3xl font-black text-foreground tracking-tight">Bitcoin</h3>
          <span className="text-lg font-medium text-muted-foreground">BTC</span>
        </div>
        {latestPrice !== null && latestPrice !== undefined && (
          <div className="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-500 mt-2 tracking-tighter">
            $
            {latestPrice.toLocaleString(undefined, {
              minimumFractionDigits: 2,
              maximumFractionDigits: 2,
            })}
          </div>
        )}
      </div>

      {/* Gradient Overlay for that "fade" effect */}
      <div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-transparent z-10 pointer-events-none" />
      <div ref={chartContainerRef} className="w-full h-full" />
    </div>
  );
};
