import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { usePortfolioValue } from './usePortfolioValue';
import { getQuote } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import { mockActiveQuote, mockPortfolio } from '@/test/mocks';

vi.mock('@/lib/api', () => ({
  getQuote: vi.fn(),
}));

describe('usePortfolioValue', () => {
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

  it('returns zero value if portfolio is empty or undefined', () => {
    const { result } = renderHook(() => usePortfolioValue(undefined), { wrapper });
    expect(result.current.totalValue).toBe(0);
    expect(result.current.isLoading).toBe(false);
  });

  it('aggregates quotes and calculates total portfolio value correctly', async () => {
    (getQuote as unknown as ReturnType<typeof vi.fn>).mockImplementation(({ symbol }) => {
      if (symbol === 'AAPL')
        return Promise.resolve({ ...mockActiveQuote, symbol: 'AAPL', price: 150 });
      if (symbol === 'MSFT')
        return Promise.resolve({ ...mockActiveQuote, symbol: 'MSFT', price: 300 });
      return Promise.resolve({ ...mockActiveQuote, price: 0 });
    });

    const { result } = renderHook(() => usePortfolioValue(mockPortfolio), { wrapper });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    // AAPL: 10 * 150 + MSFT: 5 * 300 = 3000
    expect(result.current.totalValue).toBe(3000);
    expect(result.current.isError).toBe(false);
  });
});
