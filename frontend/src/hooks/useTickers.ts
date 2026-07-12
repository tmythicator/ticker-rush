import { getActiveLadder, getQuote } from '@/lib/api';
import { useQuery, useQueries } from '@tanstack/react-query';
import { queryKeys } from '@/lib/queryKeys';
import { queryConfig } from '@/lib/queryConfig';
import { useMemo } from 'react';
import type { TradeSymbol } from '@/types';

export const useTickers = () => {
  const {
    data: ladder,
    isLoading: isLadderLoading,
    error: ladderError,
  } = useQuery({
    queryKey: queryKeys.ladder.active,
    queryFn: getActiveLadder,
    staleTime: Infinity,
    refetchOnWindowFocus: false,
  });

  const allowedTickers = useMemo(() => ladder?.allowed_tickers ?? [], [ladder]);

  // Fetch quotes for all allowed tickers in parallel to check open/closed status
  const quotesQueries = useQueries({
    queries: allowedTickers.map((t) => ({
      queryKey: queryKeys.quotes.detail(t.symbol as TradeSymbol),
      queryFn: () => getQuote({ symbol: t.symbol as TradeSymbol }),
      enabled: allowedTickers.length > 0,
      ...queryConfig.quotes,
    })),
  });

  const sortedTickers = useMemo(() => {
    if (allowedTickers.length === 0) return [];

    // Map each symbol to its closed status
    const closedMap: Record<string, boolean> = {};
    quotesQueries.forEach((q, index) => {
      const symbol = allowedTickers[index].symbol;
      if (q.data) {
        closedMap[symbol] = q.data.is_closed;
      } else {
        closedMap[symbol] = false;
      }
    });

    // Sort: open markets (is_closed == false) first, closed markets (is_closed == true) last
    return [...allowedTickers].sort((a, b) => {
      const aClosed = closedMap[a.symbol];
      const bClosed = closedMap[b.symbol];
      if (aClosed === bClosed) return 0;
      return aClosed ? 1 : -1;
    });
  }, [allowedTickers, quotesQueries]);

  const isLoading = isLadderLoading || quotesQueries.some((q) => q.isLoading);
  const error = ladderError || quotesQueries.find((q) => q.error)?.error || null;

  return {
    data: ladder?.allowed_tickers ? sortedTickers : undefined,
    isLoading,
    error,
  };
};
