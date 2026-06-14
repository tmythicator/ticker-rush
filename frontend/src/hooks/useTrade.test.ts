import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTrade } from './useTrade';
import { buyStock, sellStock } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import { TradeAction } from '@/types';
import { mockUserParticipating } from '@/test/mocks';

vi.mock('@/lib/api', () => ({
  buyStock: vi.fn(),
  sellStock: vi.fn(),
}));

describe('useTrade', () => {
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

  it('submits buyStock mutation and calls options.onSuccess', async () => {
    const successSpy = vi.fn();
    (buyStock as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);

    const { result } = renderHook(() => useTrade({ symbol: 'AAPL', onSuccess: successSpy }), {
      wrapper,
    });

    act(() => {
      result.current.executeTrade({ action: TradeAction.BUY, quantity: 10 });
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(buyStock).toHaveBeenCalledWith({ symbol: 'AAPL', quantity: 10 });
    expect(successSpy).toHaveBeenCalled();
  });

  it('submits sellStock mutation', async () => {
    (sellStock as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);

    const { result } = renderHook(() => useTrade({ symbol: 'AAPL' }), { wrapper });

    act(() => {
      result.current.executeTrade({ action: TradeAction.SELL, quantity: 5 });
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(sellStock).toHaveBeenCalledWith({ symbol: 'AAPL', quantity: 5 });
  });

  it('throws error for negative or zero quantities', async () => {
    const { result } = renderHook(() => useTrade({ symbol: 'AAPL' }), { wrapper });

    act(() => {
      result.current.executeTrade({ action: TradeAction.BUY, quantity: 0 });
    });

    await waitFor(() => {
      expect(result.current.error).toBe('Quantity must be positive');
    });
  });
});
