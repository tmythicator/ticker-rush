import { useChart } from '@/hooks/useChart';
import { type Quote, type TradeSymbol } from '@/types';
import { useRef } from 'react';
import { ChartSymbolIndicator } from './ChartSymbolIndicator';
import { ChartSymbolPicker } from './ChartSymbolPicker';

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
      <div className="flex h-full w-full items-center justify-center rounded-lg bg-muted/10 text-muted-foreground">
        No active ticker selected. Check configuration.
      </div>
    );
  }

  return (
    <div className="group relative h-full w-full">
      <div className="absolute left-4 top-4 z-20 flex flex-col gap-2">
        <div className="flex items-center gap-1 rounded-xl border border-border/60 bg-background/90 p-1.5 shadow-lg backdrop-blur-md transition-all hover:scale-[1.02] hover:shadow-xl">
          <ChartSymbolPicker symbol={symbol} onSymbolChange={onSymbolChange} />
          <div className="mx-1 h-6 w-px bg-border"></div>
          <ChartSymbolIndicator quote={quote} isLoading={isLoading} isError={isError} />
        </div>
      </div>
      <div className="group relative h-full w-full">
        <div ref={chartContainerRef} className="h-[500px] w-full" />
      </div>
    </div>
  );
};
