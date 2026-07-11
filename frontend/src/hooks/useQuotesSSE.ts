import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useState } from 'react';
import { type Quote } from '@/types';
import { queryKeys } from '@/lib/queryKeys';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';

export const useQuotesSSE = (symbol: string | null) => {
  const [error, setError] = useState<Event | null>(null);
  const queryClient = useQueryClient();

  // Initial hydration + subscriber
  const { data: quote = null } = useQuoteQuery(symbol);

  // Check event sequence timestamps and update if needed
  const updateCache = useCallback(
    (newData: Quote) => {
      queryClient.setQueryData(queryKeys.quotes.detail(newData.symbol), (oldData: Quote | null) => {
        if (!oldData || newData.timestamp > oldData.timestamp) {
          return newData;
        }
        return oldData;
      });
    },
    [queryClient],
  );

  // SSE worker
  useEffect(() => {
    if (!symbol) return;

    const url = `${import.meta.env.VITE_API_URL}/v1/quotes/events?symbol=${symbol}`;
    const eventSource = new EventSource(url);

    eventSource.onopen = (event) => {
      console.log('SSE: Connection Opened', event);
      setError(null);
    };

    eventSource.onerror = (errorEvent: Event) => {
      console.error('SSE: Connection Error', errorEvent);
      setError(errorEvent);
    };

    eventSource.addEventListener('quote', (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data) as Quote;
        updateCache(data);
      } catch (err) {
        console.error('SSE: Parse Error', err);
      }
    });

    // Cleanup handler for page transition, refresh or suspend
    const cleanup = () => {
      eventSource.close();
    };

    // Listen to both beforeunload (legacy) and pagehide (modern)
    window.addEventListener('beforeunload', cleanup);
    window.addEventListener('pagehide', cleanup);

    return () => {
      window.removeEventListener('beforeunload', cleanup);
      window.removeEventListener('pagehide', cleanup);
      cleanup();
    };
  }, [updateCache, symbol]);

  return { quote, error };
};
