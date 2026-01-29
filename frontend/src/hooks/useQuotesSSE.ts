import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useState } from 'react';
import { type Quote } from '../lib/api';
import { QUERY_KEY_QUOTE } from '../lib/queryKeys';
import { useQuoteQuery } from './useQuoteQuery';

export const useQuotesSSE = (symbol: string) => {
  const [error, setError] = useState<Event | null>(null);

  const queryClient = useQueryClient();

  // Initial hydration + subscriber
  const { data: quote = null } = useQuoteQuery(symbol);

  // State management
  const setQuote = useCallback(
    (newData: Quote) => {
      queryClient.setQueryData(QUERY_KEY_QUOTE(newData.symbol), (oldData: Quote | null) => {
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
    const url = `${import.meta.env.VITE_API_URL}/quotes/events?symbol=${symbol}`;
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
        setQuote(data);
      } catch (err) {
        console.error('SSE: Parse Error', err);
      }
    });

    return () => {
      eventSource.close();
    };
  }, [setQuote, symbol]);

  return { quote, error };
};
