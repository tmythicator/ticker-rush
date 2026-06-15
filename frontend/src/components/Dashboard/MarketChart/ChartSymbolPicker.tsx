import { IconChevronDown } from '@/components/icons/CustomIcons';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { useTickers } from '@/hooks/useTickers';
import { isTradeSymbol, type TickerSource, type TradeSymbol } from '@/types';

interface ChartSymbolPickerProps {
  symbol: TradeSymbol | null;
  onSymbolChange: (symbol: TradeSymbol) => void;
}

export const ChartSymbolPicker = ({ symbol, onSymbolChange }: ChartSymbolPickerProps) => {
  const { data: config } = useTickers();
  const tickers = config || [];

  const handleSymbolChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newSymbol = e.target.value;
    if (isTradeSymbol(newSymbol, tickers)) {
      onSymbolChange(newSymbol);
    }
  };

  const tickerInfo = tickers.find((t) => t.symbol === symbol);
  const source = (tickerInfo?.source ?? 'Finnhub') as TickerSource;

  return (
    <div className="group/select relative">
      <div className="flex items-center gap-2">
        {symbol && <SourceBadge source={source} />}
        <select
          value={symbol || ''}
          onChange={handleSymbolChange}
          disabled={tickers.length === 0}
          className="cursor-pointer appearance-none bg-transparent py-1.5 pl-1 pr-6 text-lg font-bold tracking-tight text-foreground transition-colors hover:text-primary focus:outline-none disabled:cursor-not-allowed disabled:text-muted-foreground"
        >
          {tickers.length === 0 ? (
            <option value="">No assets available</option>
          ) : (
            tickers.map((t) => (
              <option
                key={t.symbol}
                value={t.symbol}
                className="bg-popover text-popover-foreground"
              >
                {t.symbol.toUpperCase()}
              </option>
            ))
          )}
        </select>
        <IconChevronDown className="pointer-events-none absolute right-0 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground transition-colors group-hover/select:text-primary" />
      </div>
    </div>
  );
};
