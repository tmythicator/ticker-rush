import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTradePanel } from './useTradePanel';
import { useAuth } from '@/hooks/useAuth';
import { useTrade } from '@/hooks/useTrade';
import { mockActiveQuote, mockUserParticipating } from '@/test/mocks';
import { TradeAction } from '@/types';

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

vi.mock('@/hooks/useTrade', () => ({
  useTrade: vi.fn(),
}));

describe('useTradePanel', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('calculates correct values and delegates execution on handleTrade', () => {
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: mockUserParticipating,
    });

    const executeTradeSpy = vi.fn();
    (useTrade as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      executeTrade: executeTradeSpy,
      isLoading: false,
      error: null,
    });

    const { result } = renderHook(() => useTradePanel({ quote: mockActiveQuote }));

    // Initial state
    expect(result.current.form.quantity).toBe('');
    expect(result.current.asset.price).toBe(150);
    expect(result.current.asset.positionQuantity).toBe(10);
    expect(result.current.asset.buyingPower).toBe(5000);
    expect(result.current.estCost).toBe(0);

    // Set quantity
    act(() => {
      result.current.form.setQuantity('10');
    });

    expect(result.current.form.quantity).toBe('10');
    expect(result.current.estCost).toBe(1500); // 10 * 150

    // Trigger trade submission
    act(() => {
      result.current.handleTrade(TradeAction.BUY);
    });

    expect(executeTradeSpy).toHaveBeenCalledWith({
      action: TradeAction.BUY,
      quantity: 10,
    });
  });

  it('rounds quantity to 8 decimal places on submission', () => {
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: mockUserParticipating,
    });

    const executeTradeSpy = vi.fn();
    (useTrade as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      executeTrade: executeTradeSpy,
      isLoading: false,
      error: null,
    });

    const { result } = renderHook(() => useTradePanel({ quote: mockActiveQuote }));

    act(() => {
      result.current.form.setQuantity('1.000000018');
    });
    act(() => {
      result.current.handleTrade(TradeAction.BUY);
    });
    expect(result.current.form.error).toBeNull();
    expect(executeTradeSpy).toHaveBeenCalledWith({
      action: TradeAction.BUY,
      quantity: 1.00000002,
    });
  });
});
