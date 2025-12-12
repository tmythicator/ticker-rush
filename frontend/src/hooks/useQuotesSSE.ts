import { useState, useEffect } from 'react';
import { fetchQuote, type Quote } from '../lib/api';

export const useQuotesSSE = (symbol: string) => {
    const [quote, setQuote] = useState<Quote | null>(null);
    const [error, setError] = useState<Event | null>(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            console.error("SSE: No auth token found");
            return;
        }

        // Hydrate initial quote
        fetchQuote(symbol, token).then(initialQuote => {
            setQuote(initialQuote);
        }).catch(e => console.error("Initial fetch failed:", e));

        const url = `${import.meta.env.VITE_API_URL}/quotes/events?symbol=${symbol}&token=${token}`;
        const eventSource = new EventSource(url);

        eventSource.onopen = () => {
            setError(null);
        };

        eventSource.onerror = (e) => {
            console.error("SSE: Connection Error", e);
            setError(e);
            eventSource.close();
        };

        eventSource.addEventListener('quote', (event: MessageEvent) => {
            try {
                const data = JSON.parse(event.data) as Quote;
                setQuote(data);
            } catch (err) {
                console.error("SSE: Parse Error", err);
            }
        });

        return () => {
            eventSource.close();
        };
    }, [symbol]);

    return { quote, error };
};
