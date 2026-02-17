import { useRef } from 'react';
import { type Quote } from '@/types';
import { ChartBody } from './ChartBody';
import { type TradeSymbol } from '@/types';
import { useChart } from '@/hooks/useChart';
import { usePriceColor } from '@/hooks/usePriceColor';

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

  const priceColor = usePriceColor(quote?.price);

  if (!symbol) {
    return (
      <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-muted/10 rounded-lg">
        No active ticker selected. Check configuration.
      </div>
    );
  }

  return (
    <div className="w-full h-full relative group">
      <ChartBody
        symbol={symbol}
        onSymbolChange={onSymbolChange}
        price={quote?.price}
        isClosed={quote?.is_closed}
        priceColor={priceColor}
        isLoading={isLoading}
        isError={isError}
      />
      <div className="w-full h-full relative group">
        <div ref={chartContainerRef} className="w-full h-[500px]" />
      </div>
    </div>
  );
};
