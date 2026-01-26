import { useState, useEffect } from 'react';
import { fetchQuote, type Quote } from '../lib/api';

export const useQuotesSSE = (symbol: string) => {
  const [quote, setQuote] = useState<Quote | null>(null);
  const [error, setError] = useState<Event | null>(null);

  useEffect(() => {
    // Hydrate initial quote
    fetchQuote(symbol)
      .then((initialQuote) => {
        setQuote(initialQuote);
      })
      .catch((e) => console.error('Initial fetch failed:', e));

    const url = `${import.meta.env.VITE_API_URL}/quotes/events?symbol=${symbol}`;
    const eventSource = new EventSource(url);

    eventSource.onopen = () => {
      setError(null);
    };

    eventSource.onerror = (e) => {
      console.error('SSE: Connection Error', e);
      setError(e);
      eventSource.close();
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
  }, [symbol]);

  return { quote, error };
};
