interface ChartSymbolIndicatorProps {
  price: number | undefined;
  priceColor: string;
  isLoading: boolean;
  isError: boolean;
}

export const ChartSymbolIndicator = ({
  price,
  priceColor,
  isLoading,
  isError,
}: ChartSymbolIndicatorProps) => {
  return (
    <div className="px-2 flex flex-col items-end min-w-[80px]">
      {isLoading ? (
        <div className="h-5 w-16 bg-slate-200 animate-pulse rounded"></div>
      ) : isError ? (
        <span className="text-xs font-bold text-red-500">OFFLINE</span>
      ) : (
        <>
          <span
            className={`text-lg font-mono font-bold leading-none ${priceColor} transition-colors duration-300`}
          >
            {price ? `$${price.toFixed(2)}` : '—'}
          </span>
          <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
            Live
          </span>
        </>
      )}
    </div>
  );
};
