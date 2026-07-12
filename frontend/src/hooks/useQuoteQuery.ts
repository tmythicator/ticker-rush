import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '@/lib/queryKeys';
import { getQuote } from '@/lib/api';
import { queryConfig } from '@/lib/queryConfig';
import type { TradeSymbol } from '@/types';

export const useQuoteQuery = (symbol: string | null) => {
  return useQuery({
    queryKey: queryKeys.quotes.detail(symbol as TradeSymbol),
    queryFn: () => getQuote({ symbol: symbol as TradeSymbol }),
    enabled: !!symbol, // Don't fetch if symbol is empty
    ...queryConfig.quotes,
  });
};
