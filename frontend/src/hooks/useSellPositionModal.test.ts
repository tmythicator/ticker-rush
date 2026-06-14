import { vi, describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useSellPositionModal } from './useSellPositionModal';
import { useQuotesSSE } from '@/hooks/useQuotesSSE';
import { useTrade } from '@/hooks/useTrade';
import { useTickers } from '@/hooks/useTickers';
import { TradeAction } from '@/types';
import { mockActiveQuote } from '@/test/mocks';

vi.mock('@/hooks/useQuotesSSE', () => ({
  useQuotesSSE: vi.fn(),
}));

vi.mock('@/hooks/useTrade', () => ({
  useTrade: vi.fn(),
}));

vi.mock('@/hooks/useTickers', () => ({
  useTickers: vi.fn(),
}));

describe('useSellPositionModal', () => {
  let executeTradeSpy: ReturnType<typeof vi.fn>;
  let onCloseSpy: () => void;
  let onSuccessSpy: () => void;

  beforeEach(() => {
    vi.clearAllMocks();
    executeTradeSpy = vi.fn();
    onCloseSpy = vi.fn();
    onSuccessSpy = vi.fn();

    (useQuotesSSE as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      quote: mockActiveQuote,
      error: null,
    });

    (useTrade as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      executeTrade: executeTradeSpy,
      isLoading: false,
    });

    (useTickers as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: [{ symbol: 'AAPL' }],
    });
  });

  it('calculates portfolio total value and formats ticker name', () => {
    const { result } = renderHook(() =>
      useSellPositionModal({
        isOpen: true,
        symbol: 'aapl',
        quantity: 10,
        onClose: onCloseSpy,
        onSuccess: onSuccessSpy,
      }),
    );

    expect(result.current.displaySymbol).toBe('AAPL');
    expect(result.current.price).toBe(150);
    expect(result.current.totalValue).toBe(1500); // 10 * 150
  });

  it('triggers executeTrade when handleSellAll is called', () => {
    const { result } = renderHook(() =>
      useSellPositionModal({
        isOpen: true,
        symbol: 'AAPL',
        quantity: 5,
        onClose: onCloseSpy,
        onSuccess: onSuccessSpy,
      }),
    );

    act(() => {
      result.current.handleSellAll();
    });

    expect(executeTradeSpy).toHaveBeenCalledWith({
      action: TradeAction.SELL,
      quantity: 5,
    });
  });
});
