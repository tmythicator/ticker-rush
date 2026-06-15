import { usePriceColor } from '@/hooks/usePriceColor';
import { type Quote } from '@/types';

interface ChartSymbolIndicatorProps {
  quote: Quote | null;
  isLoading: boolean;
  isError: boolean;
}

export const ChartSymbolIndicator = ({ quote, isLoading, isError }: ChartSymbolIndicatorProps) => {
  const price = quote?.price;
  const isClosed = quote?.is_closed;
  const priceColor = usePriceColor(price);
  return (
    <div className="flex min-w-[80px] flex-col items-end px-2">
      {isLoading ? (
        <div className="h-5 w-16 animate-pulse rounded bg-muted"></div>
      ) : isError ? (
        <span className="text-xs font-bold text-destructive">OFFLINE</span>
      ) : (
        <>
          <span
            className={`font-mono text-lg font-bold leading-none ${priceColor} transition-colors duration-300`}
          >
            {price ? `$${price.toFixed(2)}` : '—'}
          </span>
          {isClosed ? (
            <span className="text-[10px] font-bold uppercase tracking-wider text-yellow-500">
              Market Closed
            </span>
          ) : (
            <span className="text-[10px] font-bold uppercase tracking-wider text-foreground/60">
              Live
            </span>
          )}
        </>
      )}
    </div>
  );
};
