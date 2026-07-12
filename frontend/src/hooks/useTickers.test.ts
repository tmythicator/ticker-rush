import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTickers } from './useTickers';
import { getActiveLadder, getQuote } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';

vi.mock('@/lib/api', () => ({
  getActiveLadder: vi.fn(),
  getQuote: vi.fn(),
}));

describe('useTickers', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });
  });

  const wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);

  it('queries active ladder and maps allowed_tickers sorted by market status', async () => {
    const mockAllowedTickers = [
      { symbol: 'AAPL', source: 'Finnhub' },
      { symbol: 'GOOGL', source: 'Finnhub' },
      { symbol: 'MSFT', source: 'Finnhub' },
    ];
    (getActiveLadder as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
      allowed_tickers: mockAllowedTickers,
    });

    (getQuote as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ symbol }) => {
      // Mock AAPL and MSFT as open (is_closed: false), GOOGL as closed (is_closed: true)
      if (symbol === 'GOOGL') {
        return Promise.resolve({ symbol, is_closed: true });
      }
      return Promise.resolve({ symbol, is_closed: false });
    });

    const { result } = renderHook(() => useTickers(), { wrapper });

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    // AAPL and MSFT (open) should be sorted first, GOOGL (closed) last
    expect(result.current.data).toEqual([
      { symbol: 'AAPL', source: 'Finnhub' },
      { symbol: 'MSFT', source: 'Finnhub' },
      { symbol: 'GOOGL', source: 'Finnhub' },
    ]);
  });
});
