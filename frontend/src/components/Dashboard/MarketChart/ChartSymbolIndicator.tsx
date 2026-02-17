interface ChartSymbolIndicatorProps {
  price: number | undefined;
  isClosed?: boolean;
  priceColor: string;
  isLoading: boolean;
  isError: boolean;
}

export const ChartSymbolIndicator = ({
  price,
  isClosed,
  priceColor,
  isLoading,
  isError,
}: ChartSymbolIndicatorProps) => {
  return (
    <div className="px-2 flex flex-col items-end min-w-[80px]">
      {isLoading ? (
        <div className="h-5 w-16 bg-muted animate-pulse rounded"></div>
      ) : isError ? (
        <span className="text-xs font-bold text-destructive">OFFLINE</span>
      ) : (
        <>
          <span
            className={`text-lg font-mono font-bold leading-none ${priceColor} transition-colors duration-300`}
          >
            {price ? `$${price.toFixed(2)}` : '—'}
          </span>
          {isClosed ? (
            <span className="text-[10px] font-bold text-yellow-500 uppercase tracking-wider">
              Market Closed
            </span>
          ) : (
            <span className="text-[10px] font-bold text-foreground/60 uppercase tracking-wider">
              Live
            </span>
          )}
        </>
      )}
    </div>
  );
};
