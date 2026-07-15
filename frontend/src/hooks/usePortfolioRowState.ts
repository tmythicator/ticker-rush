import { useMemo } from 'react';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { useTickers } from '@/hooks/useTickers';
import { formatCurrencyWithSign } from '@/lib/utils';
import type { PortfolioItem, TickerSource } from '@/types';

export const usePortfolioRowState = (item: PortfolioItem) => {
  const { data: config } = useTickers();
  const { data: quote } = useQuoteQuery(item.stock_symbol);

  const { symbol, source, isMarketClosed, isTradable } = useMemo(() => {
    const tickers = config ?? [];
    const tickerInfo = tickers.find(
      (t) => t.symbol.toUpperCase() === item.stock_symbol.toUpperCase(),
    );
    return {
      symbol: (tickerInfo?.symbol ?? item.stock_symbol).toUpperCase(),
      source: (tickerInfo?.source ?? quote?.source ?? 'Finnhub') as TickerSource,
      isMarketClosed: quote?.is_closed ?? false,
      isTradable: config ? !!tickerInfo : true,
    };
  }, [config, item.stock_symbol, quote?.source, quote?.is_closed]);

  const { marketValue, pnl, pnlStatus } = useMemo(() => {
    if (!quote) return { marketValue: null, pnl: null, pnlStatus: 'neutral' };

    const currentMarketValue = quote.price * item.quantity;
    const currentPnl = (quote.price - item.average_price) * item.quantity;

    return {
      marketValue: `$${currentMarketValue.toFixed(2)}`,
      pnl: formatCurrencyWithSign(currentPnl),
      pnlStatus: currentPnl >= 0 ? 'positive' : 'negative',
    };
  }, [quote, item.quantity, item.average_price]);

  return {
    symbol,
    source,
    isMarketClosed,
    isTradable,
    quote,
    marketValue,
    pnl,
    pnlStatus,
  };
};
