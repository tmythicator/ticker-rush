import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { usePortfolioTable } from './usePortfolioTable';
import { useNavigate } from 'react-router-dom';
import { mockPortfolio } from '@/test/mocks';

vi.mock('react-router-dom', () => ({
  useNavigate: vi.fn(),
}));

describe('usePortfolioTable', () => {
  it('maps portfolio record to a list of items', () => {
    const { result } = renderHook(() => usePortfolioTable(mockPortfolio));

    expect(result.current.items).toHaveLength(2);
    expect(result.current.items[0]).toEqual(mockPortfolio.AAPL);
    expect(result.current.items[1]).toEqual(mockPortfolio.MSFT);
    expect(result.current.isEmpty).toBe(false);
  });

  it('reports isEmpty as true for empty portfolio', () => {
    const { result } = renderHook(() => usePortfolioTable({}));

    expect(result.current.items).toHaveLength(0);
    expect(result.current.isEmpty).toBe(true);
  });

  it('triggers navigate on handleTrade', () => {
    const mockNavigate = vi.fn();
    vi.mocked(useNavigate).mockReturnValue(mockNavigate);

    const { result } = renderHook(() => usePortfolioTable(mockPortfolio));
    result.current.handleTrade('AAPL');

    expect(mockNavigate).toHaveBeenCalledWith('/trade?symbol=AAPL');
  });

  it('handles opening and closing the sell modal', () => {
    const { result } = renderHook(() => usePortfolioTable(mockPortfolio));

    // Initial state
    expect(result.current.sellModal.isOpen).toBe(false);
    expect(result.current.sellModal.item).toBeUndefined();

    // Open modal
    act(() => {
      result.current.handleSellClick(mockPortfolio.AAPL);
    });

    expect(result.current.sellModal.isOpen).toBe(true);
    expect(result.current.sellModal.item).toEqual(mockPortfolio.AAPL);

    // Close modal
    act(() => {
      result.current.handleCloseSellModal();
    });

    expect(result.current.sellModal.isOpen).toBe(false);
    expect(result.current.sellModal.item).toBeUndefined();
  });
});
