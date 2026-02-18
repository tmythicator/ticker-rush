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
      <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-muted/10 rounded-lg">
        No active ticker selected. Check configuration.
      </div>
    );
  }

  return (
    <div className="w-full h-full relative group">
      <div className="absolute top-4 left-4 z-20 flex flex-col gap-2">
        <div className="bg-background/90 backdrop-blur-md p-1.5 rounded-xl shadow-lg border border-border/60 flex items-center gap-1 transition-all hover:shadow-xl hover:scale-[1.02]">
          <ChartSymbolPicker symbol={symbol} onSymbolChange={onSymbolChange} />
          <div className="h-6 w-px bg-border mx-1"></div>
          <ChartSymbolIndicator quote={quote} isLoading={isLoading} isError={isError} />
        </div>
      </div>
      <div className="w-full h-full relative group">
        <div ref={chartContainerRef} className="w-full h-[500px]" />
      </div>
    </div>
  );
};
