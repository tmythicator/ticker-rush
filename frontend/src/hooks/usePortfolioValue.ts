import { fetchQuote } from '@/lib/api';
import { QUERY_KEY_QUOTE } from '@/lib/queryKeys';
import { type PortfolioItem } from '@/types';
import { useQueries } from '@tanstack/react-query';
import { useMemo } from 'react';

export const usePortfolioValue = (portfolio: Record<string, PortfolioItem> | undefined) => {
  const symbols = useMemo(() => Object.keys(portfolio || {}), [portfolio]);

  const results = useQueries({
    queries: symbols.map((symbol) => ({
      queryKey: QUERY_KEY_QUOTE(symbol),
      queryFn: () => fetchQuote(symbol),
      staleTime: 1000 * 30,
      retry: false,
    })),
  });

  const isLoading = results.some((result) => result.isLoading);
  const isError = results.some((result) => result.isError);

  const totalValue = useMemo(() => {
    if (!portfolio || isLoading || isError) return 0;

    return symbols.reduce((acc, symbol, index) => {
      const quote = results[index].data;
      const quantity = portfolio[symbol]?.quantity || 0;
      const price = quote?.price || 0;
      return acc + quantity * price;
    }, 0);
  }, [portfolio, symbols, results, isLoading, isError]);

  return { totalValue, isLoading, isError };
};
