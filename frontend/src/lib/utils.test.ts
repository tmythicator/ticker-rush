import { describe, it, expect } from 'vitest';
import { calculateMaxBuyQuantity, calculateInvestedCapital, formatCurrencyWithSign } from './utils';
import { mockPortfolio } from '@/test/mocks';

describe('calculateMaxBuyQuantity', () => {
  it('correctly calculates quantity for standard presets (e.g. 50%)', () => {
    // Buying power = $1000, price = $150, percentage = 50% (0.5)
    // expected: (1000 * 0.5) / 150 = 500 / 150 = 3.33333333... => formatted to 6 decimal places: '3.333333'
    const qty = calculateMaxBuyQuantity(1000, 150, 0.5);
    expect(qty).toBe('3.333333');
  });

  it('correctly calculates quantity for 25%', () => {
    // Buying power = $1000, price = $150, percentage = 25% (0.25)
    // expected: (1000 * 0.25) / 150 = 250 / 150 = 1.6666666... => '1.666667'
    const qty = calculateMaxBuyQuantity(1000, 150, 0.25);
    expect(qty).toBe('1.666667');
  });

  it('adjusts the percentage to 0.999999 for MAX (100% or greater) to leave a buffer', () => {
    // Buying power = $1000, price = $100, percentage = 1.0 (MAX)
    // expected cash used: 1000 * 0.999999 = $999.999
    // expected quantity: 999.999 / 100 = 9.99999 => formatted: '9.999990'
    const qty = calculateMaxBuyQuantity(1000, 100, 1.0);
    expect(qty).toBe('9.999990');

    // Greater than 1.0 should also scale to 0.999999
    const qtyGreater = calculateMaxBuyQuantity(1000, 100, 1.5);
    expect(qtyGreater).toBe('9.999990');
  });

  it('returns zero string for edge cases (zero/negative price or buying power)', () => {
    expect(calculateMaxBuyQuantity(0, 100, 0.5)).toBe('0.000000');
    expect(calculateMaxBuyQuantity(1000, 0, 0.5)).toBe('0.000000');
    expect(calculateMaxBuyQuantity(1000, -50, 0.5)).toBe('0.000000');
  });
});

describe('calculateInvestedCapital', () => {
  it('calculates the sum of average price * quantity correctly', () => {
    // Expected: 10 * 150 + 5 * 300 = 1500 + 1500 = 3000
    expect(calculateInvestedCapital(mockPortfolio)).toBe(3000);
  });

  it('handles null or empty portfolio gracefully', () => {
    expect(calculateInvestedCapital(null)).toBe(0);
    expect(calculateInvestedCapital(undefined)).toBe(0);
    expect(calculateInvestedCapital({})).toBe(0);
  });
});

describe('formatCurrencyWithSign', () => {
  it('formats positive numbers with a + sign', () => {
    expect(formatCurrencyWithSign(123.456)).toBe('+$123.46');
    expect(formatCurrencyWithSign(0)).toBe('+$0.00');
  });

  it('formats negative numbers with a - sign', () => {
    expect(formatCurrencyWithSign(-456.789)).toBe('-$456.79');
  });
});
