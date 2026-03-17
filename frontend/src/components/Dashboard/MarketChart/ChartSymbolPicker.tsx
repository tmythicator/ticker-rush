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
    <div className="relative group/select">
      <div className="flex items-center gap-2">
        {symbol && <SourceBadge source={source} />}
        <select
          value={symbol || ''}
          onChange={handleSymbolChange}
          disabled={tickers.length === 0}
          className="appearance-none bg-transparent pl-1 pr-6 py-1.5 font-bold text-foreground text-lg tracking-tight focus:outline-none cursor-pointer hover:text-primary transition-colors disabled:cursor-not-allowed disabled:text-muted-foreground"
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
        <IconChevronDown className="w-4 h-4 text-muted-foreground absolute right-0 top-1/2 -translate-y-1/2 pointer-events-none group-hover/select:text-primary transition-colors" />
      </div>
    </div>
  );
};
