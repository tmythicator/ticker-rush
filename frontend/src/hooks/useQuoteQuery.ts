import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '@/lib/queryKeys';
import { getQuote } from '@/lib/api';

export const useQuoteQuery = (symbol: string | null) => {
  return useQuery({
    queryKey: queryKeys.quotes.detail(symbol!),
    queryFn: () => getQuote({ symbol: symbol! }),
    enabled: !!symbol, // Don't fetch if symbol is empty
    staleTime: 1000 * 30, // Keep quote fresh for 30 sec
  });
};
