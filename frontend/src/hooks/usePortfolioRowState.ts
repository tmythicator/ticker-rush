import { useMemo } from 'react';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { useTickers } from '@/hooks/useTickers';
import { formatCurrencyWithSign } from '@/lib/utils';
import type { PortfolioItem, TickerSource } from '@/types';

export const usePortfolioRowState = (item: PortfolioItem) => {
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

  return {
    symbol,
    source,
    isMarketClosed,
    quote,
    marketValue,
    pnl,
    pnlColorClass,
  };
};
