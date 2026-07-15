import { renderHook } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { usePriceColor } from './usePriceColor';

describe('usePriceColor', () => {
  it('returns neutral initially', () => {
    const { result } = renderHook(() => usePriceColor(100));
    expect(result.current).toBe('neutral');
  });

  it('keeps neutral if price is undefined', () => {
    const { result } = renderHook(() => usePriceColor(undefined));
    expect(result.current).toBe('neutral');
  });

  it('returns up when price increases', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 105 });
    expect(result.current).toBe('up');
  });

  it('returns down when price decreases', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 95 });
    expect(result.current).toBe('down');
  });

  it('keeps previous color status when price remains the same', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 105 });
    expect(result.current).toBe('up');

    rerender({ price: 105 });
    expect(result.current).toBe('up');
  });
});
