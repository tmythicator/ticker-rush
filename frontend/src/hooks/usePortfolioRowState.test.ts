import { renderHook } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { usePortfolioRowState } from './usePortfolioRowState';
import { useTickers } from '@/hooks/useTickers';
import { useQuoteQuery } from '@/hooks/useQuoteQuery';
import { mockPortfolioItemAAPL } from '@/test/mocks';

vi.mock('@/hooks/useTickers', () => ({
  useTickers: vi.fn(),
}));

vi.mock('@/hooks/useQuoteQuery', () => ({
  useQuoteQuery: vi.fn(),
}));

describe('usePortfolioRowState', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('computes state correctly for allowed/tradable and active assets', () => {
    vi.mocked(useTickers).mockReturnValue({
      data: [{ symbol: 'AAPL', source: 'Finnhub' }],
      isLoading: false,
      error: null,
    });
    vi.mocked(useQuoteQuery).mockReturnValue({
      data: {
        symbol: 'AAPL',
        price: 160.0,
        change: 10.0,
        change_percent: 6.67,
        timestamp: new Date().toISOString(),
        source: 'Finnhub',
        is_closed: false,
      },
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useQuoteQuery>);

    const { result } = renderHook(() => usePortfolioRowState(mockPortfolioItemAAPL));

    expect(result.current.symbol).toBe('AAPL');
    expect(result.current.source).toBe('Finnhub');
    expect(result.current.isTradable).toBe(true);
    expect(result.current.isMarketClosed).toBe(false);
    expect(result.current.marketValue).toBe('$1600.00');
    expect(result.current.pnl).toBe('+$100.00');
    expect(result.current.pnlColorClass).toBe('text-green-500');
  });

  it('computes isTradable as false if symbol is not in allowed tickers list', () => {
    vi.mocked(useTickers).mockReturnValue({
      data: [{ symbol: 'GOOG', source: 'Finnhub' }],
      isLoading: false,
      error: null,
    });
    vi.mocked(useQuoteQuery).mockReturnValue({
      data: {
        symbol: 'AAPL',
        price: 160.0,
        change: 10.0,
        change_percent: 6.67,
        timestamp: new Date().toISOString(),
        source: 'Finnhub',
        is_closed: false,
      },
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useQuoteQuery>);

    const { result } = renderHook(() => usePortfolioRowState(mockPortfolioItemAAPL));

    expect(result.current.isTradable).toBe(false);
  });

  it('computes isMarketClosed as true if quote says the market is closed', () => {
    vi.mocked(useTickers).mockReturnValue({
      data: [{ symbol: 'AAPL', source: 'Finnhub' }],
      isLoading: false,
      error: null,
    });
    vi.mocked(useQuoteQuery).mockReturnValue({
      data: {
        symbol: 'AAPL',
        price: 160.0,
        change: 10.0,
        change_percent: 6.67,
        timestamp: new Date().toISOString(),
        source: 'Finnhub',
        is_closed: true,
      },
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useQuoteQuery>);

    const { result } = renderHook(() => usePortfolioRowState(mockPortfolioItemAAPL));

    expect(result.current.isMarketClosed).toBe(true);
  });

  it('defaults isTradable to true if tickers config has not loaded yet', () => {
    vi.mocked(useTickers).mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    });
    vi.mocked(useQuoteQuery).mockReturnValue({
      data: null,
      isLoading: true,
      error: null,
    } as unknown as ReturnType<typeof useQuoteQuery>);

    const { result } = renderHook(() => usePortfolioRowState(mockPortfolioItemAAPL));

    expect(result.current.isTradable).toBe(true);
    expect(result.current.marketValue).toBeNull();
    expect(result.current.pnl).toBeNull();
    expect(result.current.pnlColorClass).toBe('text-muted-foreground');
  });
});
