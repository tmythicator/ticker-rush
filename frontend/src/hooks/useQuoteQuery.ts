import { useQuery } from '@tanstack/react-query';
import { QUERY_KEY_QUOTE } from '../lib/queryKeys';
import { fetchQuote } from '../lib/api';

export const useQuoteQuery = (symbol: string | null) => {
  return useQuery({
    queryKey: QUERY_KEY_QUOTE(symbol!),
    queryFn: () => fetchQuote(symbol!),
    enabled: !!symbol, // Don't fetch if symbol is empty
    staleTime: 1000 * 30, // Keep quote fresh for 30 sec
  });
};
