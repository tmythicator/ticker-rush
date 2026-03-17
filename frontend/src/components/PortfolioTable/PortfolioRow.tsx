import { useMemo } from 'react';
import { IconArrowRight, IconTrash } from '@/components/icons/CustomIcons';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { Button } from '@/components/ui/button';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { useTickers } from '@/hooks/useTickers';
import { cn, formatCurrencyWithSign } from '@/lib/utils';
import type { PortfolioItem, TickerSource } from '@/types';

interface PortfolioRowProps {
  item: PortfolioItem;
  onSellClick: (item: PortfolioItem) => void;
  onTradeClick: (symbol: string) => void;
  isReadOnly?: boolean;
}

export const PortfolioRow = ({
  item,
  onSellClick,
  onTradeClick,
  isReadOnly = false,
}: PortfolioRowProps) => {
  const { data: config } = useTickers();

  const { data: quote } = useQuoteQuery(item.stock_symbol);

  const { symbol, source, isMarketClosed } = useMemo(() => {
    const tickers = config ?? [];
    const tickerInfo = tickers.find((t) => t.symbol === item.stock_symbol);
    return {
      symbol: (tickerInfo?.symbol ?? item.stock_symbol).toUpperCase(),
      source: (tickerInfo?.source ?? quote?.source ?? 'Finnhub') as TickerSource,
      isMarketClosed: quote?.is_closed ?? false,
    };
  }, [config, item.stock_symbol, quote?.source, quote?.is_closed]);

  const { marketValue, pnl, pnlColorClass } = useMemo(() => {
    if (!quote) return { marketValue: null, pnl: null, pnlColorClass: 'text-muted-foreground' };

    const currentMarketValue = quote.price * item.quantity;
    const currentPnl = (quote.price - item.average_price) * item.quantity;

    return {
      marketValue: `$${currentMarketValue.toFixed(2)}`,
      pnl: formatCurrencyWithSign(currentPnl),
      pnlColorClass: currentPnl >= 0 ? 'text-green-500' : 'text-red-500',
    };
  }, [quote, item.quantity, item.average_price]);

  const loadingIndicator = <span className="animate-pulse">...</span>;

  return (
    <tr key={item.stock_symbol} className="hover:bg-muted/50 transition-colors">
      <td className="px-6 py-4 font-bold text-foreground">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
            {symbol[0] ?? '?'}
          </div>
          <div className="flex flex-col items-start text-foreground font-bold">
            <span className="text-sm font-bold">{symbol}</span>
            <SourceBadge source={source} />
          </div>
        </div>
      </td>
      <td className="px-6 py-4 text-right font-mono text-muted-foreground">{item.quantity}</td>
      <td className="px-6 py-4 text-right font-mono text-muted-foreground">
        ${item.average_price.toFixed(2)}
      </td>
      <td className="px-6 py-4 text-right font-mono text-foreground font-medium">
        {quote ? `$${quote.price.toFixed(2)}` : loadingIndicator}
      </td>
      <td className="px-6 py-4 text-right font-mono text-foreground font-bold">
        {marketValue ?? loadingIndicator}
      </td>
      <td className={cn('px-6 py-4 text-right font-mono font-bold', pnlColorClass)}>
        {pnl ?? loadingIndicator}
      </td>
      {!isReadOnly && (
        <td className="px-6 py-4 text-right">
          <div className="flex items-center justify-end gap-2">
            <Button
              variant="destructive"
              size="sm"
              onClick={() => onSellClick(item)}
              title={isMarketClosed ? 'Market Closed' : 'Sell All'}
              className="h-8 px-3 text-xs"
              disabled={isMarketClosed}
            >
              <IconTrash className="w-3 h-3 mr-1" />
              Sell All
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={() => onTradeClick(item.stock_symbol)}
              className="h-8 px-3 text-xs"
              disabled={isMarketClosed}
              title={isMarketClosed ? 'Market Closed' : 'Trade'}
            >
              Trade
              <IconArrowRight className="w-3 h-3 ml-1" />
            </Button>
          </div>
        </td>
      )}
    </tr>
  );
};
