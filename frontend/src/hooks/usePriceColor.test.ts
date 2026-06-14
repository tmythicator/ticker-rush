import { renderHook } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { usePriceColor } from './usePriceColor';

describe('usePriceColor', () => {
  it('returns default foreground color initially', () => {
    const { result } = renderHook(() => usePriceColor(100));
    expect(result.current).toBe('text-foreground');
  });

  it('keeps default color if price is undefined', () => {
    const { result } = renderHook(() => usePriceColor(undefined));
    expect(result.current).toBe('text-foreground');
  });

  it('returns green color when price increases', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 105 });
    expect(result.current).toBe('text-green-500');
  });

  it('returns red color when price decreases', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 95 });
    expect(result.current).toBe('text-red-500');
  });

  it('keeps previous color when price remains the same', () => {
    const { result, rerender } = renderHook(({ price }) => usePriceColor(price), {
      initialProps: { price: 100 },
    });

    rerender({ price: 105 });
    expect(result.current).toBe('text-green-500');

    rerender({ price: 105 });
    expect(result.current).toBe('text-green-500');
  });
});
