import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTickers } from './useTickers';
import { getActiveLadder } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';

vi.mock('@/lib/api', () => ({
  getActiveLadder: vi.fn(),
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

  it('queries active ladder and maps allowed_tickers', async () => {
    const mockAllowedTickers = ['AAPL', 'GOOGL', 'MSFT'];
    (getActiveLadder as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
      allowed_tickers: mockAllowedTickers,
    });

    const { result } = renderHook(() => useTickers(), { wrapper });

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.data).toEqual(mockAllowedTickers);
  });
});
